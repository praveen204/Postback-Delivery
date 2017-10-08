
<?php
/*This script accepts HTTP POST requests  containing raw data in JSON format and appends a postback object to a redis list for each data object in the data */
require 'Predis/Autoloader.php';
Predis\Autoloader::register();
$redis_list = 'request';
//Checking if JSON can be decoded from the raw data
if($info = json_decode(file_get_contents('php://input'), true)){	
	//Checking if endpoint is in valid URL format or not
	if(filter_var($info['endpoint']['url'], FILTER_VALIDATE_URL)){
		//Checking if endpoint method is valid or not
		if(strtoupper($info['endpoint']['method']) == 'GET' ||
		   strtoupper($info['endpoint']['method']) == 'POST'){
			try{   //Connecting to local redis database
				$redis = new Predis\Client (array ('scheme'   => 'tcp','host'     => 'localhost','port'     => 7000,'password' => 'test',));
				//Pushing postback objects to redis list(Delivery Queue) for each data object
				foreach($info['data'] as $data){
					$redis->rpush($redis_list, json_encode($info['endpoint'] + array('data' => $data)));
				}
				echo "Pushing to Redis Successful!<br>";
			}
			catch (Exception $msg){
				//Handling Redis connection issues
				echo "Connection Error: can't connect to Redis.<br>";
				echo $msg->getMessage();
			}
		}
		else{
			echo 'Invalid method given.Use either "GET" or "POST" for method parameter.<br>';
		}
	}
	else{
		echo 'Invalid URL entered.Please enter a valid URL<br>';
	}
}
else{
	echo 'Raw data in the POST request is not in valid JSON format. Please check data again.<br>';
}
?>
