/*This program acts as a delivery agent to continuously pull postback objects on port 7000 in the server*/

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
//Pbo type helps in defining the structure to store the postback objects
// Pbo - Postback object
type Pbo struct {
	Method string              `json:"method"`
	Url    string              `json:"url"`
	Data   map[string]string   `json:"data"`
}
// The below constants store configuration related information
const (
	LOG_FILE = "delivery_agent.log"
	REDIS_LIST = "request"
	MISMATCH_KEY_VALUE_URL = ""	
	DISPLAY_TRACES = false
)
// Variables declaration
var (
	argumentPattern = regexp.MustCompile("{.*?}")
	
	v_trace   *log.Logger
	v_info    *log.Logger
	v_warning *log.Logger
	v_error   *log.Logger
)


/*The below function helps in initializing various logs */

func initializeLogs(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	v_warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
	v_error =   log.New(errorHandle,   "ERROR: ",   log.Ldate|log.Lmicroseconds|log.Llongfile)
	v_trace =   log.New(traceHandle,   "TRACE: ",   log.Ldate|log.Lmicroseconds|log.Llongfile)
	v_info =    log.New(infoHandle,    "INFO: ",    log.Ldate|log.Lmicroseconds|log.Llongfile)

}

/* The below function handles response received from sending a postback object and logs it into  info file*/
func logEndpointResponseInfo(response *http.Response, postback Pbo) {
	v_info.Println("Received response from : < " + postback.Url+" >" )
	v_info.Println("Response Code:", response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	v_info.Println("Response Body:", string(body))
}


/*The below function finds all keys in the Postback object and replaces them with and replace them with values present in Postback object's data*/
func mappingUrlKeystoValues(postback *Pbo) {
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


/*The below function delivers postback object using GET method*/

func deliverForGetType(postback Pbo) {
	requestBody, _ := json.Marshal(postback.Data)
	v_trace.Println("Request Body: " + string(requestBody))
	request, err :=  http.NewRequest("GET", postback.Url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	v_trace.Println("request: " + fmt.Sprint(request))
	v_info.Println("Delivering URL: < " + postback.Url + " >   method: " + postback.Method)
	response, err := client.Do(request)
	if err != nil {
		v_warning.Println("Could not send GET request to: " + postback.Url + ">")
	} else {
		defer response.Body.Close()
		logEndpointResponseInfo(response, postback)
	}
}


/*The below function delivers postback object using POST  method*/
func deliverForPostType(postback Pbo) {
	requestBody, _ := json.Marshal(postback.Data)
	v_trace.Println("requestBody: " + string(requestBody))
	request, err :=  http.NewRequest("POST", postback.Url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	v_trace.Println("request: " + fmt.Sprint(request))
	v_info.Println("Delivering URL: < " + postback.Url + " >  method: " + postback.Method)
	response, err := client.Do(request)
	if err != nil {
		v_warning.Println("Could not send POST request to: <" + postback.Url + ">")
	} else {
		defer response.Body.Close()
		logEndpointResponseInfo(response, postback)
	}
}
/*The below function performs the processing of postback object received from Redis server*/
func processPbo(redisServer redis.Conn) {
	endpoint, err := redis.String(redisServer.Do("LPOP", REDIS_LIST))
	if err == nil && endpoint != "" {
		postback := Pbo{}
		json.Unmarshal([]byte(endpoint), &postback)
		v_trace.Println("endpoint: " + endpoint)
		v_trace.Println("postback: " + fmt.Sprint(postback))
		mappingUrlKeystoValues(&postback)
		v_trace.Println("postback.Url: " + postback.Url)
		if strings.ToUpper(postback.Method) == "GET" {
			deliverForGetType(postback)
		} else if strings.ToUpper(postback.Method) == "POST" {
			deliverForPostType(postback)
		} else {
			v_error.Println("Unsupported Postback Method.")
		}
	} else if fmt.Sprint(err) == "redigo: nil returned" {
		//There is no data yet to deliver.
	} else if err != nil {
		v_warning.Println("Redis Problem: " + fmt.Sprint(err))
	} else {
		v_warning.Println("Received empty Redis Object.")
	}
}


/*The below function is the main function from where the actual execution starts*/
func main() {
	logger, logError := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logError != nil {
		log.Fatalln("Unable to open log file: ", logError)
	}
	defer logger.Close()

	var stdout_trace io.Writer
	if DISPLAY_TRACES {
		stdout_trace = os.Stdout
	} else {
		stdout_trace = ioutil.Discard
	}
	stdout_warn := io.MultiWriter(logger, os.Stdout)
	stdout_info := io.MultiWriter(logger, os.Stdout)	
	stdout_error := io.MultiWriter(logger, os.Stderr)
	initializeLogs(stdout_trace, stdout_info, stdout_warn, stdout_error)
	
	redisServer, err := redis.Dial("tcp", ":7000")
	if err != nil {
		v_error.Fatalln(err)
	}
	  if _, err := redisServer.Do("AUTH","test"); err != nil {
		v_error.Fatalln(err)
        }	
	defer redisServer.Close()
	for {
		processPbo(redisServer)
	}
}
