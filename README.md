###### Postback-Delivery

Introduction:
Postback delivery is a service that takes HTTP POST requests and delivers HTTP requests mentioned in the body.

It composes of three components namely:

 I) Ingestion Agent
 ii) Delivery Queue
iii) Delivery Agent

Ingestion Agent:
It is the agent which accepts incoming requests and  pushes a "postback" object to Delivery Queue for each "data" object contained in accepted request.

In this project, the PHP HTTP web server which is built on Apache2 acts as an Ingestion Agent.

Delivery Queue:

The delivery queue is the entity responsible for accepting the postback objects from ingestion agent as well as sending them to delivery agent.

In this project, Redis database is used to maintain this delivery queue.

Delivery Agent:
The delivery agent is responsible for continuously pulling postback objects from Redis and delivering each of them to the http endpoint.

In this project, the delivery agent is the concurrent golang service that keeps running until someone forcibly stops it.

Prerequisites:
i) A linux system preferably Ubuntu 16.04 as this tutorial is based on it.

ii) Ensure network settings are properly configured in it and it is pingable from outside world by making needed changes in  the file of /etc/network/interfaces

Iii) Ensure that root user of the system is SSHable from outside world by changing configuration to “PermitRootLogin yes”  in /etc/ssh/sshd_config file

iv) Make sure that gcc and make packages are installed by using the following command: 
apt-get install -y make gcc

Installation Steps:

i) Install Redis using the following commands:

wget http://download.redis.io/releases/redis-4.0.2.tar.gz
tar zxvf redis-4.0.2.tar.gz
cd redis-4.0.2/
make

Also install redis-tools  in order to be able access redis command line interface. Use the following command to install reds-tools:
apt-get install -y redis-tools

ii) Install git package using the following command:

apt-get install -y git

iii) Install Apache Server using the following command:

apt-get install -y apache2-bin

iv) Install PHP 7 using the following command:

apt-get install -y php7.0-gd libapache2-mod-php7.0 php7.0-mcrypt

v) Install Go using the following commands:
wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz

vi) Set the environment variables of PATH,GOPATH by appending the following lines in $HOME/.bashrc file at the bottom:

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$(go env GOPATH)

Then run the command “bash” to update environment variables immediately

vii) Now if we try to retrieve the values, they will look like this:

echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin

echo $GOPATH
/root/go

viii) Install Predis library which helps in connecting Redis with PHP. Install it using the following commands:

apt-get install -y php-pear
pear channel-discover pear.nrk.io
pear install nrk/Predis

ix) Install Redigo library into the system so that Go can interact with Redis using this library. Install it using the following command:

go get github.com/garyburd/redigo/redis

This will create a directory of github.com in $HOME/go/src directory which internally contains redigo directory.


Changing Default port of Redis(6379 to 7000):

i) Firstly, copy the redis-server and redis-cli to the /usr/local/bin directory so that they will be available in the environment variable “PATH”

cp  $HOME/redis-4.0.2/src/redis-cli /usr/local/bin
cp  $HOME/redis-4.0.2/src/redis-server /usr/local/bin

ii) Then, create a directory where to store your Redis config files and your data as follows:
sudo mkdir /etc/redis
sudo mkdir /var/redis

iii) Copy the init script that you'll find in the Redis distribution under the utils directory into /etc/init.d directory:

sudo cp utils/redis_init_script /etc/init.d/redis_7000

iv) Edit the init script:

sudo vi /etc/init.d/redis_7000

 Set REDISPORT to 7000

REDISPORT=7000

Both the pid file path and the configuration file name depend on the port number.

v) Copy the template configuration file you'll find in the root directory of the Redis distribution into /etc/redis/ using the port number as name, for example:

sudo cp redis.conf /etc/redis/7000.conf

vi) Create a directory inside /var/redis that will work as data and working directory for this Redis instance:

sudo mkdir /var/redis/7000

vii) Edit the configuration file  /etc/redis/7000.conf as follows:

sudo vi /etc/redis/7000.conf 

Change the following in the opened file:

- Set daemonize to yes
- Set the pidfile to /var/run/redis_7000.pid
- Change the port to 7000
- Set the logfile to /var/log/redis_7000.log
- Set the dir to /var/redis/7000


viii) 
Add the following comments containing init info in the configuration file of redis_7000 just below the #!/bin/sh line:

### BEGIN INIT INFO
# Provides: redis_7000
# Required-Start: $remote_fs $syslog
# Required-Stop: $remote_fs $syslog
# Default-Start: 2 3 4 5
# Default-Stop: 0 1 6
# Short-Description: redis_7000
# Description: This file is the  service file for redis configuration
#
### END INIT INFO

ix) Finally add the new Redis init script to all the default runlevels using the following command:
sudo update-rc.d redis_7000 defaults

x)  Then restart init service as follows:

/etc/init.d/redis_7000 stop
/etc/init.d/redis_7000 start

A reboot is recommended if any issue arises in this step

Adding authentication to Redis:

A password can be added to the redis database by uncommenting the “requirepass” property in the configuration file.

In our setup, this property is uncommented in /etc/redis/7000.conf file as follows:

requirepass test

Then restart init service as follows:

/etc/init.d/redis_7000 stop
/etc/init.d/redis_7000 start


So, any client that needs to connect to redis database should provide this password.
For example, client can connect through redid-cli with password as follows:
redis-cli -p 7000 -a test


Configuration details:
- REDIS_LIST property in delivery_agent.go and redis_list variable of ingestion_agent.go are related to each other and should be of same value
- The default log file name is delivery_agent.log but can be changed by setting the property of LOG_FILE to a different value in delivery_agent.go
- Setting the property of SHOW_TRACES to true if we want to display the trace level logs in output.

Stepwise procedure to run:
i) First we need to clone the git repository to local system if not already cloned using the following command:

go get github.com/praveen204/Postback-delivery

ii) We can access the cloned repository in $HOME/go/src/github.com/praveen204 directory

iii) Copy ingest.php and printMethod.php to /var/www/html directory using the following commands:
cp ingest.php  /var/www/html
cp printMethod.php /var/www/html

iv)  Now open POSTMAN and submit a POST request of URL http://165.227.0.65/ingest.php with following data in body:
 {  
      "endpoint":{  
        "method":"GET",
        "url":"http://localhost/printMethod.php/data?title={mascot}&image={location}&foo={bar}"
      },
      "data":[  
        {  
          "mascot":"Gopher",
          "location":"https://blog.golang.org/gopher/gopher.png"
        }
      ]
    }

v) Before submitting, make sure to run the delivery_agent.go to see response details using the following command:

go run delivery_agent.go

vi) Now click submit in POSTMAN for the request. A response of “Pushing to Redis Successful!” appears in POSTMAN.

vii) Each of the post back objects generated using the request’s data  will be pushed to redis list and response will be generated. The running delivery_agent.go shows the logs of result in the console for our reference.

INFO: 2017/10/08 05:32:50.600462 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:88: Delivering URL: < http://localhost/printMethod.php/data?title=Gopher&image=https://blog.golang.org/gopher/gopher.png&foo=  >  method: GET
INFO: 2017/10/08 05:32:50.602125 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:55: Received response from: < http://localhost/printMethod.php/data?title=Gopher&image=https://blog.golang.org/gopher/gopher.png&foo= >
INFO: 2017/10/08 05:32:50.602150 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:56: Response Code: 200
INFO: 2017/10/08 05:32:50.602252 /root/go/src/github.com/praveen204/Postback-delivery/delivery_agent.go:58: Response Body: {"bar":"","location":"https://blog.golang.org/gopher/gopher.png","mascot":"Gopher"}

Sample Run:

Sample request:

(POST) http://{server_ip}/ingest.php
    (RAW POST DATA) 
    {  
      "endpoint":{  
        "method":"GET",
        "url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"
      },
      "data":[  
        {  
          "mascot":"Gopher",
          "location":"https://blog.golang.org/gopher/gopher.png"
        }
      ]
    }


Sample response:

 GET http://sample_domain_endpoint.com/data?title=Gopher&image=https://blog.golang.org/gopher/gopher.png&foo=




