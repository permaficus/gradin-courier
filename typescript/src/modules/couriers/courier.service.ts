import { NotFoundError, ValidationError } from "../../common/errors.js";
import type { CourierPayload, CourierQuery, CourierResponse } from "./courier.model.js";
import { courierPayloadSchema, parseCourierQuery, zodErrors } from "./courier.schema.js";

interface CourierStore {
  create(payload: CourierPayload): Promise<CourierResponse>;
  list(query: CourierQuery): Promise<{
    data: CourierResponse[];
    pagination: { page: number; per_page: number; total: number; total_pages: number };
  }>;
  findById(id: string): Promise<CourierResponse | null>;
  update(id: string, payload: CourierPayload): Promise<CourierResponse | null>;
  softDelete(id: string): Promise<boolean>;
}

export class CourierService {
  constructor(private readonly repository: CourierStore) {}

  async create(payload: unknown) {
    const parsed = courierPayloadSchema.safeParse(payload);
    if (!parsed.success) {
      throw new ValidationError(zodErrors(parsed.error));
    }
    return this.repository.create(parsed.data);
  }

  async list(query: Record<string, unknown>) {
    const parsed = parseCourierQuery(query);
    if (parsed.errors || !parsed.value) {
      throw new ValidationError(parsed.errors ?? {}, 400);
    }
    return this.repository.list(parsed.value);
  }

  async find(id: string) {
    assertObjectId(id);
    const courier = await this.repository.findById(id);
    if (!courier) {
      throw new NotFoundError();
    }
    return courier;
  }

  async update(id: string, payload: unknown) {
    assertObjectId(id);
    const parsed = courierPayloadSchema.safeParse(payload);
    if (!parsed.success) {
      throw new ValidationError(zodErrors(parsed.error));
    }
    const courier = await this.repository.update(id, parsed.data as CourierPayload);
    if (!courier) {
      throw new NotFoundError();
    }
    return courier;
  }

  async delete(id: string) {
    assertObjectId(id);
    const deleted = await this.repository.softDelete(id);
    if (!deleted) {
      throw new NotFoundError();
    }
  }
}

function assertObjectId(id: string) {
  if (!/^[a-f\d]{24}$/i.test(id)) {
    throw new NotFoundError();
  }
}
