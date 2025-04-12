<?php
require_once 'vendor/autoload.php';

$host = getenv('REDIS_HOST') ? getenv('REDIS_HOST') : 'localhost';
$port = getenv('REDIS_PORT') ? getenv('REDIS_PORT') : '6379';

$dsn = "redis://${host}:{$port}";
echo "Connecting to $dsn";
Resque::setBackend($dsn);

for ($x = 0; $x <= 30; $x++) {
  $num = $x % 5;
  Resque::enqueue('default', "My_Job$num", ['name' => 'Test', 'val' => $x]);
}

?>