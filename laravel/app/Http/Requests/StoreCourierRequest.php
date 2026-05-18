<?php

namespace App\Http\Requests;

use App\Enums\CourierLevel;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class StoreCourierRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'min:2', 'max:150'],
            'email' => ['nullable', 'email', 'max:150'],
            'phone' => ['nullable', 'string', 'max:30'],
            'level' => ['required', 'integer', Rule::in(CourierLevel::values())],
            'vehicle_type' => ['nullable', 'string', 'max:50'],
            'license_plate' => ['nullable', 'string', 'max:30'],
            'status' => ['nullable', Rule::in(['active', 'inactive', 'suspended'])],
            'registered_at' => ['nullable', 'date'],
        ];
    }
}
