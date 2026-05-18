<?php

namespace App\Http\Controllers;

use App\Http\Requests\StoreCourierRequest;
use App\Http\Requests\UpdateCourierRequest;
use App\Http\Resources\CourierResource;
use App\Services\CourierService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Str;
use Illuminate\Validation\ValidationException;

class CourierController
{
    public function __construct(private readonly CourierService $couriers)
    {
    }

    public function index(Request $request): JsonResponse
    {
        try {
            $page = $this->couriers->list($request->query());
        } catch (ValidationException $exception) {
            return response()->json([
                'success' => false,
                'message' => 'Validation failed',
                'errors' => $exception->errors(),
                'meta' => $this->meta($request),
            ], 400);
        }

        return response()->json([
            'success' => true,
            'message' => 'Couriers retrieved successfully',
            'data' => CourierResource::collection($page->items()),
            'pagination' => [
                'page' => $page->currentPage(),
                'per_page' => $page->perPage(),
                'total' => $page->total(),
                'total_pages' => $page->lastPage(),
            ],
            'meta' => $this->meta($request),
        ]);
    }

    public function store(StoreCourierRequest $request): JsonResponse
    {
        $courier = $this->couriers->create($request->validated());

        return response()->json([
            'success' => true,
            'message' => 'Courier created successfully',
            'data' => new CourierResource($courier),
            'meta' => $this->meta($request),
        ], 201);
    }

    public function show(Request $request, string $courier): JsonResponse
    {
        $model = $this->couriers->find($courier);
        if (! $model) {
            return $this->notFound($request);
        }

        return response()->json([
            'success' => true,
            'message' => 'Courier retrieved successfully',
            'data' => new CourierResource($model),
            'meta' => $this->meta($request),
        ]);
    }

    public function update(UpdateCourierRequest $request, string $courier): JsonResponse
    {
        $model = $this->couriers->find($courier);
        if (! $model) {
            return $this->notFound($request);
        }

        return response()->json([
            'success' => true,
            'message' => 'Courier updated successfully',
            'data' => new CourierResource($this->couriers->update($model, $request->validated())),
            'meta' => $this->meta($request),
        ]);
    }

    public function destroy(Request $request, string $courier): JsonResponse
    {
        $model = $this->couriers->find($courier);
        if (! $model) {
            return $this->notFound($request);
        }

        $this->couriers->delete($model);

        return response()->json([
            'success' => true,
            'message' => 'Courier deleted successfully',
            'data' => null,
            'meta' => $this->meta($request),
        ]);
    }

    private function notFound(Request $request): JsonResponse
    {
        return response()->json([
            'success' => false,
            'message' => 'Courier not found',
            'meta' => $this->meta($request),
        ], 404);
    }

    private function meta(Request $request): array
    {
        $requestId = $request->header('x-request-id') ?: (string) Str::uuid();
        header('x-request-id: '.$requestId);

        return [
            'request_id' => $requestId,
            'timestamp' => now('UTC')->toJSON(),
        ];
    }
}
