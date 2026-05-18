import { NotFoundError, ValidationError } from "../../common/errors.js";
import type { CourierPayload } from "./courier.model.js";
import type { CourierRepository } from "./courier.repository.js";
import { courierPayloadSchema, parseCourierQuery, zodErrors } from "./courier.schema.js";

export class CourierService {
  constructor(private readonly repository: CourierRepository) {}

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
      throw new ValidationError(parsed.errors ?? {});
    }
    return this.repository.list(parsed.value);
  }

  async find(id: string) {
    const courier = await this.repository.findById(id);
    if (!courier) {
      throw new NotFoundError();
    }
    return courier;
  }

  async update(id: string, payload: unknown) {
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
    const deleted = await this.repository.softDelete(id);
    if (!deleted) {
      throw new NotFoundError();
    }
  }
}
