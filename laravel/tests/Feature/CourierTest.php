<?php

namespace Tests\Feature;

use App\Models\Courier;
use Illuminate\Support\Carbon;
use Tests\TestCase;

class CourierTest extends TestCase
{
    protected function setUp(): void
    {
        parent::setUp();
        Courier::query()->delete();
    }

    public function test_it_can_create_courier(): void
    {
        $response = $this->postJson('/api/couriers', ['name' => 'Budiono Hadi Agung', 'level' => 2]);
        $response->assertCreated();
        $this->assertDatabaseHas('couriers', ['name' => 'Budiono Hadi Agung', 'level' => 2]);
    }

    public function test_it_can_list_couriers_with_pagination(): void
    {
        Courier::query()->create(['name' => 'Budi', 'level' => 2, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers')->assertOk()->assertJsonPath('pagination.page', 1);
    }

    public function test_it_sorts_couriers_by_name_by_default(): void
    {
        Courier::query()->create(['name' => 'Zaki', 'level' => 1, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        Courier::query()->create(['name' => 'Agung', 'level' => 1, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers')->assertJsonPath('data.0.name', 'Agung');
    }

    public function test_it_can_sort_couriers_by_registered_at(): void
    {
        Courier::query()->create(['name' => 'New', 'level' => 1, 'status' => 'active', 'registered_at' => Carbon::parse('2026-05-18'), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        Courier::query()->create(['name' => 'Old', 'level' => 1, 'status' => 'active', 'registered_at' => Carbon::parse('2026-05-17'), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers?sort=registered_at')->assertJsonPath('data.0.name', 'Old');
    }

    public function test_it_can_search_courier_by_partial_words(): void
    {
        Courier::query()->create(['name' => 'Budiono Hadi Agung', 'level' => 2, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers?search=budi+agung')->assertJsonPath('data.0.name', 'Budiono Hadi Agung');
    }

    public function test_it_can_filter_couriers_by_multiple_levels(): void
    {
        Courier::query()->create(['name' => 'Level Two', 'level' => 2, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        Courier::query()->create(['name' => 'Level Five', 'level' => 5, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers?level=2,3')->assertJsonCount(1, 'data');
    }

    public function test_it_can_filter_couriers_by_single_level(): void
    {
        Courier::query()->create(['name' => 'Level Two', 'level' => 2, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        Courier::query()->create(['name' => 'Level Three', 'level' => 3, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);

        $this->getJson('/api/couriers?level=2')
            ->assertOk()
            ->assertJsonCount(1, 'data')
            ->assertJsonPath('data.0.level', 2);
    }

    public function test_it_paginates_without_duplicate_records(): void
    {
        foreach (range(1, 12) as $number) {
            Courier::query()->create([
                'name' => sprintf('Courier %02d', $number),
                'level' => 1,
                'status' => 'active',
                'registered_at' => now(),
                'created_at' => now(),
                'updated_at' => now(),
                'deleted_at' => null,
            ]);
        }

        $firstPage = $this->getJson('/api/couriers?page=1&per_page=10')
            ->assertOk()
            ->assertJsonPath('pagination.page', 1)
            ->assertJsonPath('pagination.per_page', 10)
            ->assertJsonPath('pagination.total', 12)
            ->json('data');

        $secondPage = $this->getJson('/api/couriers?page=2&per_page=10')
            ->assertOk()
            ->assertJsonPath('pagination.page', 2)
            ->assertJsonCount(2, 'data')
            ->json('data');

        $this->assertEmpty(array_intersect(
            array_column($firstPage, 'id'),
            array_column($secondPage, 'id')
        ));
    }

    public function test_it_can_show_courier(): void
    {
        $courier = Courier::query()->create(['name' => 'Shown', 'level' => 1, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->getJson('/api/couriers/'.$courier->_id)->assertOk()->assertJsonPath('data.name', 'Shown');
    }

    public function test_it_can_update_courier(): void
    {
        $courier = Courier::query()->create(['name' => 'Before', 'level' => 1, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->putJson('/api/couriers/'.$courier->_id, ['name' => 'After', 'level' => 3])->assertOk();
        $this->assertDatabaseHas('couriers', ['name' => 'After', 'level' => 3]);
    }

    public function test_it_can_delete_courier(): void
    {
        $courier = Courier::query()->create(['name' => 'Deleted', 'level' => 1, 'status' => 'active', 'registered_at' => now(), 'created_at' => now(), 'updated_at' => now(), 'deleted_at' => null]);
        $this->deleteJson('/api/couriers/'.$courier->_id)->assertOk();
        $this->getJson('/api/couriers/'.$courier->_id)->assertNotFound();
        $this->getJson('/api/couriers')->assertJsonMissing(['name' => 'Deleted']);
    }

    public function test_it_rejects_invalid_level(): void
    {
        $this->getJson('/api/couriers?level=9')->assertBadRequest();
    }

    public function test_it_rejects_invalid_query_parameters(): void
    {
        $this->getJson('/api/couriers?page=abc')->assertBadRequest();
        $this->getJson('/api/couriers?per_page=101')->assertBadRequest();
        $this->getJson('/api/couriers?sort=deleted_at')->assertBadRequest();
    }

    public function test_it_rejects_invalid_payload(): void
    {
        $this->postJson('/api/couriers', ['name' => '', 'level' => 9])->assertUnprocessable();
        $this->postJson('/api/couriers', ['level' => 2])->assertUnprocessable();
        $this->postJson('/api/couriers', ['name' => 'Invalid Email', 'level' => 2, 'email' => 'invalid'])->assertUnprocessable();
        $this->postJson('/api/couriers', ['name' => 'Invalid Status', 'level' => 2, 'status' => 'unknown'])->assertUnprocessable();
        $this->postJson('/api/couriers', ['name' => str_repeat('a', 151), 'level' => 2])->assertUnprocessable();
    }

    public function test_it_returns_not_found_for_invalid_id(): void
    {
        $this->getJson('/api/couriers/not-a-valid-object-id')->assertNotFound();
    }

    public function test_it_reuses_request_id_header(): void
    {
        $this->withHeader('x-request-id', 'test-request-id')
            ->getJson('/api/couriers')
            ->assertOk()
            ->assertHeader('x-request-id', 'test-request-id')
            ->assertJsonPath('meta.request_id', 'test-request-id');
    }
}
