<?php

namespace App\Providers;

use App\Models\Courier;
use Illuminate\Support\ServiceProvider;

class MongoIndexServiceProvider extends ServiceProvider
{
    public function boot(): void
    {
        if ($this->app->runningInConsole() && ! $this->app->runningUnitTests()) {
            return;
        }

        $collection = Courier::raw();
        $collection->createIndex(['name' => 1], ['name' => 'couriers_name_idx']);
        $collection->createIndex(['level' => 1], ['name' => 'couriers_level_idx']);
        $collection->createIndex(['registered_at' => 1], ['name' => 'couriers_registered_at_idx']);
        $collection->createIndex(['email' => 1], ['name' => 'couriers_email_unique_sparse_idx', 'unique' => true, 'sparse' => true]);
    }
}
