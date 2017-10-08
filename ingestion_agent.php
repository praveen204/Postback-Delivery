<?php
require 'Predis/Autoloader.php';
Predis\Autoloader::register();
$redis_list = 'request';
//tries to decode JSON from the raw post data. If it can not, echos an error
if($info = json_decode(file_get_contents('php://input'), true)){	
	//varifies that the provided endpoint url is in valid form. Else echos an error
	if(filter_var($info['endpoint']['url'], FILTER_VALIDATE_URL)){
		//varifies the endpoint method is supported. Else echos an error
		if(strtoupper($info['endpoint']['method']) == 'GET' ||
		   strtoupper($info['endpoint']['method']) == 'POST'){
			try{   //makes a coonection to the local redis database
				$redis = new Predis\Client (array ('scheme'   => 'tcp','host'     => 'localhost','port'     => 7000,'password' => 'test',));
				//makes separate "postback" objects for each received "data" object
				foreach($info['data'] as $data){
					$redis->rpush($redis_list, json_encode($info['endpoint'] + array('data' => $data)));
				}
				echo "Pushing to Redis Successful!<br>";
			}
			catch (Exception $msg){
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
