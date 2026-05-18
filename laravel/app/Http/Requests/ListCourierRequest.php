<?php

namespace App\Http\Requests;

use App\Support\ApiResponse;
use Illuminate\Contracts\Validation\Validator;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Http\Exceptions\HttpResponseException;
use Illuminate\Validation\Rule;

class ListCourierRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'page' => ['nullable', 'integer', 'min:1'],
            'per_page' => ['nullable', 'integer', 'min:1', 'max:100'],
            'sort' => ['nullable', Rule::in(['name', '-name', 'registered_at', '-registered_at', 'created_at', '-created_at'])],
            'search' => ['nullable', 'string'],
            'level' => ['nullable', 'string', function (string $attribute, mixed $value, \Closure $fail): void {
                if (trim((string) $value) === '') {
                    return;
                }

                foreach (explode(',', (string) $value) as $level) {
                    $integerLevel = filter_var(trim($level), FILTER_VALIDATE_INT);
                    if ($integerLevel === false || $integerLevel < 1 || $integerLevel > 5) {
                        $fail('The level query must contain only levels 1 to 5.');

                        return;
                    }
                }
            }],
        ];
    }

    public function messages(): array
    {
        return [
            'page.integer' => 'The page query is invalid.',
            'page.min' => 'The page query is invalid.',
            'per_page.integer' => 'The per_page query is invalid.',
            'per_page.min' => 'The per_page query is invalid.',
            'per_page.max' => 'The per_page query is invalid.',
            'sort.in' => 'The sort field is invalid.',
        ];
    }

    protected function failedValidation(Validator $validator): void
    {
        throw new HttpResponseException(ApiResponse::validation($this, $validator->errors()->toArray(), 400));
    }
}
