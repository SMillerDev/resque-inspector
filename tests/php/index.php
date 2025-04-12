<?php
require_once 'vendor/autoload.php';

$port = getenv('REDIS_PORT') ?? '6379';

Resque::setBackend("localhost:{$port}");

$args = array(
          'name' => 'Chris'
        );

for ($x = 0; $x <= 30; $x++) {
  $num = $x % 5;
  $val = Resque::enqueue('default', "My_Job$num", $args);
}


var_dump(Resque::queues());
?>