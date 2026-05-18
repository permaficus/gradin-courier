<?php

namespace App\Services;

use App\Models\Courier;
use App\Repositories\CourierRepository;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Support\Carbon;
use Illuminate\Validation\ValidationException;

class CourierService
{
    public function __construct(private readonly CourierRepository $couriers) {}

    public function list(array $query): LengthAwarePaginator
    {
        return $this->couriers->paginate($this->parseQuery($query));
    }

    public function create(array $data): Courier
    {
        $now = Carbon::now('UTC');

        return $this->couriers->create(array_merge($data, [
            'status' => $data['status'] ?? 'active',
            'registered_at' => isset($data['registered_at']) ? Carbon::parse($data['registered_at']) : $now,
            'created_at' => $now,
            'updated_at' => $now,
            'deleted_at' => null,
        ]));
    }

    public function find(string $id): ?Courier
    {
        return $this->couriers->findActive($id);
    }

    public function update(Courier $courier, array $data): Courier
    {
        $data['status'] = $data['status'] ?? 'active';
        $data['registered_at'] = isset($data['registered_at']) ? Carbon::parse($data['registered_at']) : $courier->registered_at;
        $data['updated_at'] = Carbon::now('UTC');

        return $this->couriers->update($courier, $data);
    }

    public function updateById(string $id, array $data): ?Courier
    {
        $courier = $this->find($id);
        if (! $courier) {
            return null;
        }

        return $this->update($courier, $data);
    }

    public function delete(Courier $courier): void
    {
        $this->couriers->softDelete($courier);
    }

    public function deleteById(string $id): bool
    {
        $courier = $this->find($id);
        if (! $courier) {
            return false;
        }

        $this->delete($courier);

        return true;
    }

    private function parseQuery(array $query): array
    {
        $page = max((int) ($query['page'] ?? 1), 1);
        $perPage = max(min((int) ($query['per_page'] ?? 10), 100), 1);
        $sort = $query['sort'] ?? 'name';
        $allowedSorts = ['name', '-name', 'registered_at', '-registered_at', 'created_at', '-created_at'];
        if (! in_array($sort, $allowedSorts, true)) {
            throw ValidationException::withMessages(['sort' => 'The sort field is invalid.']);
        }

        $levels = [];
        if (($query['level'] ?? '') !== '') {
            foreach (explode(',', (string) $query['level']) as $level) {
                $value = filter_var(trim($level), FILTER_VALIDATE_INT);
                if ($value === false || $value < 1 || $value > 5) {
                    throw ValidationException::withMessages(['level' => 'The level query must contain only levels 1 to 5.']);
                }
                $levels[] = $value;
            }
        }

        $sortDirection = str_starts_with($sort, '-') ? 'desc' : 'asc';
        $sortField = ltrim($sort, '-');

        return [
            'page' => $page,
            'per_page' => $perPage,
            'sort_field' => $sortField,
            'sort_direction' => $sortDirection,
            'levels' => $levels,
            'search_words' => preg_split('/\s+/', trim((string) ($query['search'] ?? '')), -1, PREG_SPLIT_NO_EMPTY),
        ];
    }
}
