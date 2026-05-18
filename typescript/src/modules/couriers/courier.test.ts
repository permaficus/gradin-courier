import { describe, expect, it } from "vitest";
import Fastify from "fastify";
import { NotFoundError } from "../../common/errors.js";
import { buildApp } from "../../app.js";
import { errorHandler } from "../../common/error-handler.js";
import type { CourierPayload, CourierQuery, CourierResponse } from "./courier.model.js";
import { CourierController } from "./courier.controller.js";
import { courierPayloadSchema, parseCourierQuery } from "./courier.schema.js";
import { CourierService } from "./courier.service.js";

describe("courier query validation", () => {
  it("parses multiple levels", () => {
    const parsed = parseCourierQuery({ level: "2,3" });
    expect(parsed.errors).toBeUndefined();
    expect(parsed.value?.levels).toEqual([2, 3]);
  });

  it("rejects invalid levels", () => {
    const parsed = parseCourierQuery({ level: "2,9" });
    expect(parsed.errors?.level).toBeDefined();
  });

  it("rejects arbitrary sort columns", () => {
    const parsed = parseCourierQuery({ sort: "deleted_at" });
    expect(parsed.errors?.sort).toBeDefined();
  });

  it("uses default pagination", () => {
    const parsed = parseCourierQuery({});
    expect(parsed.value?.page).toBe(1);
    expect(parsed.value?.per_page).toBe(10);
  });
});

describe("courier payload validation", () => {
  it("accepts valid create payload", () => {
    const parsed = courierPayloadSchema.safeParse({ name: "Budiono Hadi Agung", level: 2, email: "budi@example.com" });
    expect(parsed.success).toBe(true);
  });

  it("rejects invalid payload", () => {
    const parsed = courierPayloadSchema.safeParse({ name: "", level: 9, email: "invalid" });
    expect(parsed.success).toBe(false);
  });
});

describe("courier service", () => {
  it("treats invalid object ids as not found before calling the repository", async () => {
    const service = new CourierService({
      create: async () => {
        throw new Error("not used");
      },
      list: async () => {
        throw new Error("not used");
      },
      findById: async () => {
        throw new Error("repository should not be called");
      },
      update: async () => {
        throw new Error("not used");
      },
      softDelete: async () => {
        throw new Error("not used");
      }
    });

    await expect(service.find("not-a-valid-object-id")).rejects.toBeInstanceOf(NotFoundError);
  });
});

describe("courier routes", () => {
  it("returns 400 for invalid list query validation", async () => {
    const app = await buildApp();
    const response = await app.inject({ method: "GET", url: "/api/couriers?level=9" });
    await app.close();

    expect(response.statusCode).toBe(400);
    expect(response.json()).toMatchObject({
      success: false,
      message: "Validation failed",
      errors: {
        level: ["The level query must contain only levels 1 to 5."]
      }
    });
  });
});

describe("courier route integration", () => {
  it("supports create, list, show, update, delete, and soft-delete exclusion", async () => {
    const { app } = await buildTestCourierApp();

    const createResponse = await app.inject({
      method: "POST",
      url: "/api/couriers",
      payload: { name: "Budiono Hadi Agung", email: "budi-ts-integration@example.com", level: 2 }
    });
    expect(createResponse.statusCode).toBe(201);
    const created = createResponse.json().data;

    const listResponse = await app.inject({ method: "GET", url: "/api/couriers" });
    expect(listResponse.statusCode).toBe(200);
    expect(listResponse.json().data).toHaveLength(1);

    const showResponse = await app.inject({ method: "GET", url: `/api/couriers/${created.id}` });
    expect(showResponse.statusCode).toBe(200);
    expect(showResponse.json().data.name).toBe("Budiono Hadi Agung");

    const updateResponse = await app.inject({
      method: "PUT",
      url: `/api/couriers/${created.id}`,
      payload: { name: "Budi Updated", level: 3, status: "active" }
    });
    expect(updateResponse.statusCode).toBe(200);
    expect(updateResponse.json().data.level).toBe(3);

    const deleteResponse = await app.inject({ method: "DELETE", url: `/api/couriers/${created.id}` });
    expect(deleteResponse.statusCode).toBe(200);

    const deletedShowResponse = await app.inject({ method: "GET", url: `/api/couriers/${created.id}` });
    expect(deletedShowResponse.statusCode).toBe(404);

    const deletedListResponse = await app.inject({ method: "GET", url: "/api/couriers" });
    expect(deletedListResponse.json().data).toHaveLength(0);

    await app.close();
  });

  it("supports list sorting, search, filtering, and pagination", async () => {
    const { app, store } = await buildTestCourierApp();
    store.seed(
      courierFixture("Budiono Hadi Agung", 2, "2026-05-18"),
      courierFixture("Budi Santoso", 2, "2026-05-17"),
      courierFixture("Agung Prasetyo", 3, "2026-05-16"),
      courierFixture("Rudi Hartono", 4, "2026-05-15")
    );

    await expectNames(app, "/api/couriers", ["Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung", "Rudi Hartono"]);
    await expectNames(app, "/api/couriers?sort=registered_at", ["Rudi Hartono", "Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung"]);
    await expectNames(app, "/api/couriers?search=budi+agung", ["Budiono Hadi Agung"]);
    await expectNames(app, "/api/couriers?level=2", ["Budi Santoso", "Budiono Hadi Agung"]);
    await expectNames(app, "/api/couriers?level=2,3", ["Agung Prasetyo", "Budi Santoso", "Budiono Hadi Agung"]);

    const firstPage = await app.inject({ method: "GET", url: "/api/couriers?page=1&per_page=2" });
    const secondPage = await app.inject({ method: "GET", url: "/api/couriers?page=2&per_page=2" });
    expect(firstPage.statusCode).toBe(200);
    expect(secondPage.statusCode).toBe(200);
    expect(firstPage.json().pagination).toMatchObject({ page: 1, per_page: 2, total: 4, total_pages: 2 });
    expect(new Set([...ids(firstPage.json()), ...ids(secondPage.json())]).size).toBe(4);

    await app.close();
  });

  it("returns documented validation and invalid-id statuses", async () => {
    const { app } = await buildTestCourierApp();

    for (const url of ["/api/couriers?level=9", "/api/couriers?page=abc", "/api/couriers?per_page=101", "/api/couriers?sort=deleted_at"]) {
      const response = await app.inject({ method: "GET", url });
      expect(response.statusCode).toBe(400);
    }

    const invalidBodies = [
      { level: 2 },
      { name: "", level: 2 },
      { name: "Invalid Level", level: 9 },
      { name: "Invalid Email", level: 2, email: "invalid" },
      { name: "Invalid Status", level: 2, status: "unknown" },
      { name: "a".repeat(151), level: 2 }
    ];
    for (const payload of invalidBodies) {
      const response = await app.inject({ method: "POST", url: "/api/couriers", payload });
      expect(response.statusCode).toBe(422);
    }

    const invalidIdResponse = await app.inject({ method: "GET", url: "/api/couriers/not-a-valid-object-id" });
    expect(invalidIdResponse.statusCode).toBe(404);

    await app.close();
  });
});

class InMemoryCourierStore {
  private couriers = new Map<string, CourierResponse>();

  seed(...couriers: CourierResponse[]) {
    for (const courier of couriers) {
      this.couriers.set(courier.id, courier);
    }
  }

  async create(payload: CourierPayload): Promise<CourierResponse> {
    const now = new Date();
    const courier: CourierResponse = {
      id: objectId(),
      name: payload.name,
      email: payload.email ?? null,
      phone: payload.phone ?? null,
      level: payload.level,
      vehicle_type: payload.vehicle_type ?? null,
      license_plate: payload.license_plate ?? null,
      status: payload.status ?? "active",
      registered_at: payload.registered_at ? new Date(payload.registered_at) : now,
      created_at: now,
      updated_at: now,
      deleted_at: null
    };
    this.couriers.set(courier.id, courier);
    return courier;
  }

  async list(query: CourierQuery) {
    const filtered = [...this.couriers.values()]
      .filter((courier) => courier.deleted_at === null)
      .filter((courier) => query.levels.length === 0 || query.levels.includes(courier.level))
      .filter((courier) => searchWordsMatch(courier, query.search))
      .sort((left, right) => compareCouriers(left, right, query.sort));
    const start = (query.page - 1) * query.per_page;
    const data = filtered.slice(start, start + query.per_page);
    return {
      data,
      pagination: {
        page: query.page,
        per_page: query.per_page,
        total: filtered.length,
        total_pages: Math.ceil(filtered.length / query.per_page)
      }
    };
  }

  async findById(id: string): Promise<CourierResponse | null> {
    const courier = this.couriers.get(id);
    return courier && courier.deleted_at === null ? courier : null;
  }

  async update(id: string, payload: CourierPayload): Promise<CourierResponse | null> {
    const current = await this.findById(id);
    if (!current) {
      return null;
    }
    const updated = {
      ...current,
      name: payload.name,
      email: payload.email ?? null,
      phone: payload.phone ?? null,
      level: payload.level,
      vehicle_type: payload.vehicle_type ?? null,
      license_plate: payload.license_plate ?? null,
      status: payload.status ?? "active",
      registered_at: payload.registered_at ? new Date(payload.registered_at) : current.registered_at,
      updated_at: new Date()
    };
    this.couriers.set(id, updated);
    return updated;
  }

  async softDelete(id: string): Promise<boolean> {
    const current = await this.findById(id);
    if (!current) {
      return false;
    }
    this.couriers.set(id, { ...current, deleted_at: new Date(), updated_at: new Date() });
    return true;
  }
}

async function buildTestCourierApp() {
  const app = Fastify({ logger: false });
  const store = new InMemoryCourierStore();
  const controller = new CourierController(new CourierService(store));
  app.setErrorHandler(errorHandler);
  app.get("/api/couriers", controller.index);
  app.post("/api/couriers", controller.store);
  app.get<{ Params: { id: string } }>("/api/couriers/:id", controller.show);
  app.put<{ Params: { id: string } }>("/api/couriers/:id", controller.update);
  app.delete<{ Params: { id: string } }>("/api/couriers/:id", controller.destroy);
  return { app, store };
}

function courierFixture(name: string, level: number, registeredAt: string): CourierResponse {
  const now = new Date();
  return {
    id: objectId(),
    name,
    email: null,
    phone: null,
    level,
    vehicle_type: null,
    license_plate: null,
    status: "active",
    registered_at: new Date(`${registeredAt}T00:00:00.000Z`),
    created_at: now,
    updated_at: now,
    deleted_at: null
  };
}

function objectId() {
  return Array.from({ length: 24 }, () => Math.floor(Math.random() * 16).toString(16)).join("");
}

function searchWordsMatch(courier: CourierResponse, search?: string) {
  const words = search?.trim().split(/\s+/).filter(Boolean) ?? [];
  return words.every((word) => courier.name.toLowerCase().includes(word.toLowerCase()));
}

function compareCouriers(left: CourierResponse, right: CourierResponse, sort: CourierQuery["sort"]) {
  const direction = sort.startsWith("-") ? -1 : 1;
  const field = sort.replace("-", "");
  if (field === "registered_at") {
    return direction * (left.registered_at.getTime() - right.registered_at.getTime());
  }
  if (field === "created_at") {
    return direction * (left.created_at.getTime() - right.created_at.getTime());
  }
  return direction * left.name.localeCompare(right.name);
}

async function expectNames(app: Awaited<ReturnType<typeof buildTestCourierApp>>["app"], url: string, expected: string[]) {
  const response = await app.inject({ method: "GET", url });
  expect(response.statusCode).toBe(200);
  expect(response.json().data.map((courier: CourierResponse) => courier.name)).toEqual(expected);
}

function ids(responseBody: { data: CourierResponse[] }) {
  return responseBody.data.map((courier) => courier.id);
}
