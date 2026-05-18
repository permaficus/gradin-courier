<?php

namespace App\Models;

use MongoDB\Laravel\Eloquent\Model;

class Courier extends Model
{
    protected $connection = 'mongodb';

    protected $collection = 'couriers';

    protected $fillable = [
        'name',
        'email',
        'phone',
        'level',
        'vehicle_type',
        'license_plate',
        'status',
        'registered_at',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

    protected $casts = [
        'level' => 'integer',
        'registered_at' => 'datetime',
        'created_at' => 'datetime',
        'updated_at' => 'datetime',
        'deleted_at' => 'datetime',
    ];
}
