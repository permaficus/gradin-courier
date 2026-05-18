<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class CourierResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => (string) $this->_id,
            'name' => $this->name,
            'email' => $this->email,
            'phone' => $this->phone,
            'level' => $this->level,
            'vehicle_type' => $this->vehicle_type,
            'license_plate' => $this->license_plate,
            'status' => $this->status,
            'registered_at' => optional($this->registered_at)->toJSON(),
            'created_at' => optional($this->created_at)->toJSON(),
            'updated_at' => optional($this->updated_at)->toJSON(),
            'deleted_at' => optional($this->deleted_at)->toJSON(),
        ];
    }
}
