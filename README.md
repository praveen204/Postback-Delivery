{\rtf1\ansi\ansicpg1252\cocoartf1504\cocoasubrtf830
{\fonttbl\f0\fswiss\fcharset0 Helvetica;\f1\fmodern\fcharset0 Courier;}
{\colortbl;\red255\green255\blue255;\red0\green0\blue0;\red38\green38\blue38;\red255\green255\blue255;
\red254\green254\blue254;}
{\*\expandedcolortbl;;\csgenericrgb\c0\c0\c0;\cssrgb\c20000\c20000\c20000;\cssrgb\c100000\c100000\c100000;
\cssrgb\c99608\c99608\c99608;}
\margl1440\margr1440\vieww25400\viewh16000\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\qc\partightenfactor0

\f0\b\fs48 \cf2 \ul \ulc2 Postback-Delivery
\b0 \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0
\cf2 \

\b Introduction:
\b0 \

\fs28 \ulnone Postback delivery is a service that takes HTTP POST requests and delivers HTTP requests mentioned in the body.\
\
It composes of three components namely:\
\
 I) Ingestion Agent\
 ii) Delivery Queue\
iii) Delivery Agent\

\fs48 \ul \

\b Ingestion Agent:
\b0 \

\fs28 \ulnone It is the agent which accepts incoming requests and  pushes a "postback" object to Delivery Queue for each "data" object contained in accepted request.\
\
In this project, the PHP HTTP web server which is built on Apache2 acts as an Ingestion Agent.\
\

\b\fs48 \ul Delivery Queue
\fs28 \ulnone :
\b0 \
\
The delivery queue is the entity responsible for accepting the postback objects from ingestion agent as well as sending them to delivery agent.\
\
In this project, Redis database is used to maintain this delivery queue.\

\fs48 \ul \

\b Delivery Agent:
\b0 \

\fs28 \ulnone The delivery agent is responsible for continuously pulling postback objects from Redis and delivering each of them to the http endpoint.\
\
In this project, the delivery agent is the concurrent golang service that keeps running until someone forcibly stops it.\

\fs48 \ul \

\b Prerequisites:
\b0 \

\fs28 \ulnone i) A linux system preferably Ubuntu 16.04 as this tutorial is based on it.\
\
ii) Ensure network settings are properly configured in it and it is pingable from outside world by making needed changes in  the file of /etc/network/interfaces\
\
Iii) Ensure that root user of the system is SSHable from outside world by changing configuration to \'93PermitRootLogin yes\'94  in /etc/ssh/sshd_config file\
\
iv) Make sure that gcc and make packages are installed by using the following command: \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\qc\partightenfactor0

\b \cf2 apt-get install -y make gcc
\b0\fs48 \ul \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 \
Installation Steps:\
\

\b0\fs28 \ulnone i) Install Redis using the following commands:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 \
wget http://download.redis.io/releases/redis-4.0.2.tar.gz\
tar zxvf redis-4.0.2.tar.gz\
cd redis-4.0.2/\
make\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 Also install redis-tools  in order to be able access redis command line interface. Use the following command to install reds-tools:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 apt-get install -y redis-tools
\b0 \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\fs48 \cf2 \ul \

\fs28 \ulnone ii) Install git package using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 apt-get install -y git
\b0\fs48 \ul \
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\fs28 \cf2 \ulnone iii) Install Apache Server using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 apt-get install -y apache2-bin\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 iv) Install PHP 7 using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 apt-get install -y php7.0-gd libapache2-mod-php7.0 php7.0-mcrypt\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 v) Install Go using the following commands:
\fs48 \ul \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b\fs28 \cf2 \ulnone wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
\fs48 \ul \

\fs28 \ulnone tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
\b0\fs48 \ul \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\fs28 \cf2 \ulnone \
vi) Set the environment variables of PATH,GOPATH by appending the following lines in $HOME/.bashrc file at the bottom:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 export PATH=$PATH:/usr/local/go/bin\
export GOPATH=$(go env GOPATH)\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 Then run the command \'93bash\'94 to update environment variables immediately\
\
vii) Now if we try to retrieve the values, they will look like this:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 \
echo $PATH
\b0 \
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin\
\

\b echo $GOPATH
\b0 \
/root/go\
\
viii) Install Predis library which helps in connecting Redis with PHP. Install it using the following commands:\
\

\b apt-get install -y php-pear\
pear channel-discover pear.nrk.io\
pear install nrk/Predis\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 ix) Install Redigo library into the system so that Go can interact with Redis using this library. Install it using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 go get github.com/garyburd/redigo/redis\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 \
This will create a directory of github.com in $HOME/go/src directory which internally contains redigo directory.\
\
\

\b\fs48 \ul Changing Default port of Redis(6379 to 7000):\
\

\b0\fs28 \ulnone i) Firstly, copy the redis-server and redis-cli to the /usr/local/bin directory so that they will be available in the environment variable \'93PATH\'94\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 cp  $HOME/redis-4.0.2/src/redis-cli /usr/local/bin\
cp  $HOME/redis-4.0.2/src/redis-server /usr/local/bin
\b0 \
\
ii) Then, create \cf3 \cb4 \expnd0\expndtw0\kerning0
a directory where to store your Redis config files and your data as follows:\
\pard\pardeftab720\sl380\partightenfactor0

\b \cf3 \cb5 sudo mkdir /etc/redis\
sudo mkdir /var/redis\
\

\b0 iii) \cb4 Copy the init script that you'll find in the Redis distribution under the\'a0utils\'a0directory into /etc/init.d directory:
\fs32 \

\f1 \

\f0\b\fs28 \cb5 sudo cp utils/redis_init_script /etc/init.d/redis_7000\
\

\b0 iv) \cb4 Edit the init script:
\fs32 \

\b\fs28 \cb5 \
sudo vi /etc/init.d/redis_7000\

\b0 \
 Set REDISPORT to 7000\
\cb4 \

\b REDISPORT=7000\
\
\pard\pardeftab720\sl360\partightenfactor0

\b0 \cf3 Both the pid file path and the configuration file name depend on the port number.
\b \
\pard\pardeftab720\sl380\partightenfactor0
\cf3 \

\b0 v) Copy the template configuration file you'll find in the root directory of the Redis distribution into /etc/redis/ using the port number as name, for example:\
\

\b \cb5 sudo cp redis.conf /etc/redis/7000.conf
\f1\b0\fs32 \cb4 \
\

\f0\fs28 vi) Create a directory inside /var/redis that will work as data and working directory for this Redis instance:\
\pard\pardeftab720\sl380\partightenfactor0

\b\fs32 \cf3 \cb5 \
\pard\pardeftab720\sl380\partightenfactor0

\fs28 \cf3 sudo mkdir /var/redis/7000
\fs32 \cb4 \

\b0\fs28 \
vii) Edit the configuration file  
\b \cb5 /etc/redis/7000.conf 
\b0 \cb4 as follows:\
\

\b sudo vi \cb5 /etc/redis/7000.conf 
\b0 \cb4 \
\
Change the following in the opened file:\
\
- Set\'a0daemonize\'a0to yes\
- Set the\'a0pidfile\'a0to\'a0/var/run/redis_7000.pid\
- Change the port to 7000\
- Set the\'a0logfile\'a0to\'a0/var/log/redis_7000.log\
- Set the\'a0dir\'a0to /var/redis/7000\
\pard\pardeftab720\sl380\partightenfactor0

\fs32 \cf3 \cb1 \
\pard\pardeftab720\sl380\partightenfactor0

\fs28 \cf3 \cb4 \
viii) \
Add the following comments containing init info in the configuration file of redis_7000 just below the 
\b #!/bin/sh 
\b0 line:\
\

\b ### BEGIN INIT INFO\
# Provides: redis_7000\
# Required-Start: $remote_fs $syslog\
# Required-Stop: $remote_fs $syslog\
# Default-Start: 2 3 4 5\
# Default-Stop: 0 1 6\
# Short-Description: redis_7000\
# Description: This file is the  service file for redis configuration\
#\
### END INIT INFO
\b0 \
\
ix) Finally add the new Redis init script to all the default runlevels using the following command:\

\b \cb5 sudo update-rc.d redis_7000 defaults\
\

\b0 x)  \cf2 \cb1 \kerning1\expnd0\expndtw0 Then restart init service as follows:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0
\cf2 \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 /etc/init.d/redis_7000 stop\
/etc/init.d/redis_7000 start\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\fs48 \cf2 \ul \

\b0\fs28 \ulnone A reboot is recommended if any issue arises in this step
\b\fs48 \ul \
\
Adding authentication to Redis:\
\

\b0\fs28 \ulnone A password can be added to the redis database by uncommenting the \'93requirepass\'94 property in the configuration file.\
\
In our setup, this property is uncommented in /etc/redis/7000.conf file as follows:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 \
requirepass test
\b0 \
\
Then restart init service as follows:\
\

\b /etc/init.d/redis_7000 stop\
/etc/init.d/redis_7000 start\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 \
\
So, any client that needs to connect to redis database should provide this password.\
For example, client can connect through redid-cli with password as follows:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 redis-cli -p 7000 -a test
\b0 \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\fs48 \cf2 \ul \
\

\b Configuration details:\

\b0\fs28 \ulnone - REDIS_LIST property in delivery_agent.go and redis_list variable of ingestion_agent.go are related to each other and should be of same value\
- The default log file name is delivery_agent.log but can be changed by setting the property of LOG_FILE to a different value in delivery_agent.go\
- Setting the property of SHOW_TRACES to true if we want to display the trace level logs in output.\

\fs48 \ul \
\pard\pardeftab720\sl360\partightenfactor0

\b \cf2 Stepwise procedure to run:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0\fs28 \cf2 \ulnone i) First we need to clone the git repository to local system if not already cloned using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 go get github.com/praveen204/Postback-delivery\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 ii) We can access the cloned repository in $HOME/go/src/github.com/praveen204 directory\
\
iii) Copy ingest.php and printMethod.php to /var/www/html directory using the following commands:
\b \
cp \cf2 ingest.php  /var/www/html\
cp printMethod.php /var/www/html\cf2 \
\

\b0 iv)  Now open POSTMAN and submit a POST request of URL http://165.227.0.65/ingest.php with following data in body:\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2  \{  \
      "endpoint":\{  \
        "method":"GET",\
        "url":"http://localhost/printMethod.php/data?title=\{mascot\}&image=\{location\}&foo=\{bar\}"\
      \},\
      "data":[  \
        \{  \
          "mascot":"Gopher",\
          "location":"https://blog.golang.org/gopher/gopher.png"\
        \}\
      ]\
    \}\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0
\cf2 \ul \

\b0 \ulnone v) Before submitting, make sure to run the delivery_agent.go to see response details using the following command:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 go run delivery_agent.go
\b0 \
\
vi) Now click submit in POSTMAN for the request. A response of 
\b \'93Pushing to Redis Successful!\'94 
\b0 appears in POSTMAN.\
\
vii) Each of the post back objects generated using the request\'92s data  will be pushed to redis list and response will be generated. The running delivery_agent.go shows the logs of result in the console for our reference.\
\
INFO: 2017/10/08 05:32:50.600462 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:88: 
\b Delivering URL: < http://localhost/printMethod.php/data?title=Gopher&image=https://blog.golang.org/gopher/gopher.png&foo=  >  method: GET
\b0 \
INFO: 2017/10/08 05:32:50.602125 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:55: 
\b Received response from: < http://localhost/printMethod.php/data?title=Gopher&image=https://blog.golang.org/gopher/gopher.png&foo= >
\b0 \
INFO: 2017/10/08 05:32:50.602150 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:56: 
\b Response Code: 200
\b0 \
INFO: 2017/10/08 05:32:50.602252 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:58: 
\b Response Body: \{"bar":"","location":"https://blog.golang.org/gopher/gopher.png","mascot":"Gopher"\}\ul \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b0 \cf2 \

\b\fs48 Sample Run:\
\

\b0\fs28 \ulnone Sample request:\
\
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 (POST) http://\{server_ip\}/ingest.php\
    (RAW POST DATA) \
    \{  \
      "endpoint":\{  \
        "method":"GET",\
        "url":"http://sample_domain_endpoint.com/data?title=\{mascot\}&image=\{location\}&foo=\{bar\}"\
      \},\
      "data":[  \
        \{  \
          "mascot":"Gopher",\
          "location":"https://blog.golang.org/gopher/gopher.png"\
        \}\
      ]\
    \}
\b0 \
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\b \cf2 \ul \
\

\b0 \ulnone Sample response:\
\
 
\b GET http://sample_domain_endpoint.com/data?title=Gopher&image=https\cf2 ://blog.golang.org/gopher/gopher.png&foo=
\b0\fs48 \cf2 \
\
\ul \
\
}