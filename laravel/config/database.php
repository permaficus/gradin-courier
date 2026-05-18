<?php

return [
    'default' => 'mongodb',

    'connections' => [
        'mongodb' => [
            'driver' => 'mongodb',
            'dsn' => env('MONGODB_URI', 'mongodb://mongodb:mongodb12345@localhost:27019/?authSource=admin'),
            'database' => env('MONGODB_DATABASE', 'gradin-courier'),
        ],
    ],
];
