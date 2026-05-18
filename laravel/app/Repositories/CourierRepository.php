<?php

namespace App\Repositories;

use App\Models\Courier;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;
use Illuminate\Support\Carbon;

class CourierRepository
{
    public function paginate(array $filters): LengthAwarePaginator
    {
        $query = Courier::query()->whereNull('deleted_at');

        if ($filters['levels'] !== []) {
            $query->whereIn('level', $filters['levels']);
        }

        foreach ($filters['search_words'] as $word) {
            $query->where('name', 'regex', new \MongoDB\BSON\Regex(preg_quote($word), 'i'));
        }

        $query->orderBy($filters['sort_field'], $filters['sort_direction']);

        return $query->paginate($filters['per_page'], ['*'], 'page', $filters['page']);
    }

    public function create(array $data): Courier
    {
        return Courier::query()->create($data);
    }

    public function findActive(string $id): ?Courier
    {
        return Courier::query()->where('_id', $id)->whereNull('deleted_at')->first();
    }

    public function update(Courier $courier, array $data): Courier
    {
        $courier->fill($data);
        $courier->save();

        return $courier->refresh();
    }

    public function softDelete(Courier $courier): void
    {
        $courier->forceFill([
            'deleted_at' => Carbon::now('UTC'),
            'updated_at' => Carbon::now('UTC'),
        ])->save();
    }
}
