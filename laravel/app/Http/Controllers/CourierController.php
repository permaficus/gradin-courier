<?php

namespace App\Http\Controllers;

use App\Http\Requests\ListCourierRequest;
use App\Http\Requests\StoreCourierRequest;
use App\Http\Requests\UpdateCourierRequest;
use App\Http\Resources\CourierResource;
use App\Services\CourierService;
use App\Support\ApiResponse;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CourierController
{
    public function __construct(private readonly CourierService $couriers) {}

    public function index(ListCourierRequest $request): JsonResponse
    {
        $page = $this->couriers->list($request->query());

        return ApiResponse::list(
            $request,
            'Couriers retrieved successfully',
            CourierResource::collection($page->items()),
            [
                'page' => $page->currentPage(),
                'per_page' => $page->perPage(),
                'total' => $page->total(),
                'total_pages' => $page->lastPage(),
            ]
        );
    }

    public function store(StoreCourierRequest $request): JsonResponse
    {
        $courier = $this->couriers->create($request->validated());

        return ApiResponse::success($request, 'Courier created successfully', new CourierResource($courier), 201);
    }

    public function show(Request $request, string $courier): JsonResponse
    {
        $model = $this->couriers->find($courier);
        if (! $model) {
            return ApiResponse::notFound($request);
        }

        return ApiResponse::success($request, 'Courier retrieved successfully', new CourierResource($model));
    }

    public function update(UpdateCourierRequest $request, string $courier): JsonResponse
    {
        $model = $this->couriers->updateById($courier, $request->validated());
        if (! $model) {
            return ApiResponse::notFound($request);
        }

        return ApiResponse::success($request, 'Courier updated successfully', new CourierResource($model));
    }

    public function destroy(Request $request, string $courier): JsonResponse
    {
        if (! $this->couriers->deleteById($courier)) {
            return ApiResponse::notFound($request);
        }

        return ApiResponse::success($request, 'Courier deleted successfully', null);
    }
}
