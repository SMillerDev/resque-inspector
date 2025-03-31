<?php
require_once 'vendor/autoload.php';

$port = getenv('REDIS_PORT') ?? '6379';

Resque::setBackend("localhost:{$port}");

$args = array(
          'name' => 'Chris'
        );
$val = Resque::enqueue('default', 'My_Job', $args);
var_dump($val);

var_dump(Resque::queues());
?>