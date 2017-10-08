package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"bytes"		
	"os"
	"net/http"
	"io"	
	"regexp"	
	"encoding/json"
	"strings"
	"log"
	"io/ioutil"
)
//Used to store the JSON Postback object input from Redis

type Pbo struct {
	Method string              `json:"method"`
	Url    string              `json:"url"`
	Data   map[string]string   `json:"data"`
}
const (
	LOG_FILE = "delivery_agent.log"
	REDIS_LIST = "request"
	MISMATCH_KEY_VALUE_URL = ""	
	DISPLAY_TRACES = false
)

var (
	argumentPattern = regexp.MustCompile("{.*?}")
	
	V_trace   *log.Logger
	V_info    *log.Logger
	V_warning *log.Logger
	V_error   *log.Logger
)



func initializeLogs(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	V_warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
	V_error =   log.New(errorHandle,   "ERROR: ",   log.Ldate|log.Lmicroseconds|log.Llongfile)
	V_trace =   log.New(traceHandle,   "TRACE: ",   log.Ldate|log.Lmicroseconds|log.Llongfile)
	V_info =    log.New(infoHandle,    "INFO: ",    log.Ldate|log.Lmicroseconds|log.Llongfile)

}

func matchUrlKeysToValues(postback *Pbo) {
	matchingIndexes := argumentPattern.FindStringIndex(postback.Url)
	for matchingIndexes != nil {
		patternMatch := argumentPattern.FindString(postback.Url)
		matchString := patternMatch[1:(len(patternMatch) - 1)]
		replaceString, keyHasValue := postback.Data[matchString]
		if !keyHasValue {
			replaceString = MISMATCH_KEY_VALUE_URL
			postback.Data[matchString] = MISMATCH_KEY_VALUE_URL
		}
		postback.Url = postback.Url[:matchingIndexes[0]] + replaceString + postback.Url[matchingIndexes[1]:]
		matchingIndexes = argumentPattern.FindStringIndex(postback.Url)
	}
}

func logEndpointResponseInfo(response *http.Response, postback Pbo) {
	V_info.Println("Received response from: <" + postback.Url + ">")
	V_info.Println("Response Code:", response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	V_info.Println("Response Body:", string(body))
}

func deliverForGetType(postback Pbo) {
	requestBody, _ := json.Marshal(postback.Data)
	V_trace.Println("Request Body: " + string(requestBody))
	request, err :=  http.NewRequest("GET", postback.Url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	V_trace.Println("request: " + fmt.Sprint(request))
	V_info.Println("Delivering. url: <" + postback.Url + "> method: " + postback.Method)
	response, err := client.Do(request)
	if err != nil {
		V_warning.Println("Could not send GET request to: <" + postback.Url + ">")
	} else {
		defer response.Body.Close()
		logEndpointResponseInfo(response, postback)
	}
}

func deliverForPostType(postback Pbo) {
	requestBody, _ := json.Marshal(postback.Data)
	V_trace.Println("requestBody: " + string(requestBody))
	request, err :=  http.NewRequest("POST", postback.Url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	V_trace.Println("request: " + fmt.Sprint(request))
	V_info.Println("Delivering. url: <" + postback.Url + "> method: " + postback.Method)
	response, err := client.Do(request)
	if err != nil {
		V_warning.Println("Could not send POST request to: <" + postback.Url + ">")
	} else {
		defer response.Body.Close()
		logEndpointResponseInfo(response, postback)
	}
}

func processPbo(redisServer redis.Conn) {
	endpoint, err := redis.String(redisServer.Do("LPOP", REDIS_LIST))
	if err == nil && endpoint != "" {
		postback := Pbo{}
		json.Unmarshal([]byte(endpoint), &postback)
		V_trace.Println("endpoint: " + endpoint)
		V_trace.Println("postback: " + fmt.Sprint(postback))
		matchUrlKeysToValues(&postback)
		V_trace.Println("postback.Url: " + postback.Url)
		if strings.ToUpper(postback.Method) == "GET" {
			deliverForGetType(postback)
		} else if strings.ToUpper(postback.Method) == "POST" {
			deliverForPostType(postback)
		} else {
			V_error.Println("Unsupported Postback Method.")
		}
	} else if fmt.Sprint(err) == "redigo: nil returned" {
		//There is no data yet to deliver.
	} else if err != nil {
		V_warning.Println("Redis Problem: " + fmt.Sprint(err))
	} else {
		V_warning.Println("Received empty Redis Object.")
	}
}


func main() {
	logger, logError := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logError != nil {
		log.Fatalln("Unable to open log file: ", logError)
	}
	defer logger.Close()

	var traceOutput io.Writer
	if DISPLAY_TRACES {
		traceOutput = os.Stdout
	} else {
		traceOutput = ioutil.Discard
	}
	stdout_warn := io.MultiWriter(logger, os.Stdout)
	stdout_info := io.MultiWriter(logger, os.Stdout)	
	stdout_error := io.MultiWriter(logger, os.Stderr)
	initializeLogs(traceOutput, stdout_info, stdout_warn, stdout_error)
	
	redisServer, err := redis.Dial("tcp", ":7000")
	if err != nil {
		V_error.Fatalln(err)
	}
	  if _, err := redisServer.Do("AUTH","test"); err != nil {
		V_error.Fatalln(err)
        }	
	defer redisServer.Close()
	for {
		processPbo(redisServer)
	}
}
