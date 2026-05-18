import type { PrismaClient } from "@prisma/client";
import type { CourierPayload, CourierQuery, CourierResponse } from "./courier.model.js";

export class CourierRepository {
  constructor(private readonly database: PrismaClient) {}

  async create(payload: CourierPayload): Promise<CourierResponse> {
    const now = new Date();
    const courier = await this.database.courier.create({
      data: {
        name: payload.name,
        email: payload.email ?? null,
        phone: payload.phone ?? null,
        level: payload.level,
        vehicleType: payload.vehicle_type ?? null,
        licensePlate: payload.license_plate ?? null,
        status: payload.status ?? "active",
        registeredAt: payload.registered_at ? new Date(payload.registered_at) : now,
        createdAt: now,
        updatedAt: now,
        deletedAt: null
      }
    });
    return toResponse(courier);
  }

  async list(query: CourierQuery) {
    const where = whereFromQuery(query);
    const [total, couriers] = await Promise.all([
      this.database.courier.count({ where }),
      this.database.courier.findMany({
        where,
        orderBy: orderByFromSort(query.sort),
        skip: (query.page - 1) * query.per_page,
        take: query.per_page
      })
    ]);
    return {
      data: couriers.map(toResponse),
      pagination: {
        page: query.page,
        per_page: query.per_page,
        total,
        total_pages: Math.ceil(total / query.per_page)
      }
    };
  }

  async findById(id: string): Promise<CourierResponse | null> {
    const courier = await this.database.courier.findFirst({ where: { id, deletedAt: null } });
    return courier ? toResponse(courier) : null;
  }

  async update(id: string, payload: CourierPayload): Promise<CourierResponse | null> {
    const current = await this.database.courier.findFirst({ where: { id, deletedAt: null } });
    if (!current) {
      return null;
    }
    const courier = await this.database.courier.update({
      where: { id },
      data: {
        name: payload.name,
        email: payload.email ?? null,
        phone: payload.phone ?? null,
        level: payload.level,
        vehicleType: payload.vehicle_type ?? null,
        licensePlate: payload.license_plate ?? null,
        status: payload.status ?? "active",
        registeredAt: payload.registered_at ? new Date(payload.registered_at) : current.registeredAt
      }
    });
    return toResponse(courier);
  }

  async softDelete(id: string): Promise<boolean> {
    const current = await this.database.courier.findFirst({ where: { id, deletedAt: null } });
    if (!current) {
      return false;
    }
    await this.database.courier.update({ where: { id }, data: { deletedAt: new Date() } });
    return true;
  }
}

function whereFromQuery(query: CourierQuery) {
  const where: Record<string, unknown> = { deletedAt: null };
  if (query.levels.length > 0) {
    where.level = { in: query.levels };
  }
  const words = query.search?.trim().split(/\s+/).filter(Boolean) ?? [];
  if (words.length > 0) {
    where.AND = words.map((word) => ({ name: { contains: word, mode: "insensitive" } }));
  }
  return where;
}

function orderByFromSort(sort: CourierQuery["sort"]) {
  const direction = sort.startsWith("-") ? "desc" : "asc";
  const field = sort.replace("-", "");
  const fields: Record<string, string> = {
    name: "name",
    registered_at: "registeredAt",
    created_at: "createdAt"
  };
  return { [fields[field] ?? "name"]: direction };
}

function toResponse(courier: {
  id: string;
  name: string;
  email: string | null;
  phone: string | null;
  level: number;
  vehicleType: string | null;
  licensePlate: string | null;
  status: string;
  registeredAt: Date;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
}): CourierResponse {
  return {
    id: courier.id,
    name: courier.name,
    email: courier.email,
    phone: courier.phone,
    level: courier.level,
    vehicle_type: courier.vehicleType,
    license_plate: courier.licensePlate,
    status: courier.status,
    registered_at: courier.registeredAt,
    created_at: courier.createdAt,
    updated_at: courier.updatedAt,
    deleted_at: courier.deletedAt
  };
}
