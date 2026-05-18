<?php

namespace App\Support;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Str;

class ApiResponse
{
    public static function success(Request $request, string $message, mixed $data, int $status = 200): JsonResponse
    {
        return self::json($request, [
            'success' => true,
            'message' => $message,
            'data' => $data,
            'meta' => self::meta($request),
        ], $status);
    }

    public static function list(Request $request, string $message, mixed $data, array $pagination): JsonResponse
    {
        return self::json($request, [
            'success' => true,
            'message' => $message,
            'data' => $data,
            'pagination' => $pagination,
            'meta' => self::meta($request),
        ]);
    }

    public static function validation(Request $request, array $errors, int $status): JsonResponse
    {
        return self::json($request, [
            'success' => false,
            'message' => 'Validation failed',
            'errors' => $errors,
            'meta' => self::meta($request),
        ], $status);
    }

    public static function notFound(Request $request, string $message = 'Courier not found'): JsonResponse
    {
        return self::json($request, [
            'success' => false,
            'message' => $message,
            'meta' => self::meta($request),
        ], 404);
    }

    private static function json(Request $request, array $payload, int $status = 200): JsonResponse
    {
        return response()->json($payload, $status)->header('x-request-id', self::requestId($request));
    }

    private static function meta(Request $request): array
    {
        return [
            'request_id' => self::requestId($request),
            'timestamp' => now('UTC')->toJSON(),
        ];
    }

    private static function requestId(Request $request): string
    {
        if ($request->attributes->has('request_id')) {
            return (string) $request->attributes->get('request_id');
        }

        $requestId = $request->header('x-request-id') ?: (string) Str::uuid();
        $request->attributes->set('request_id', $requestId);

        return $requestId;
    }
}
