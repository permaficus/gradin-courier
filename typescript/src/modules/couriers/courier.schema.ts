import { z } from "zod";
import type { CourierQuery } from "./courier.model.js";

export const courierPayloadSchema = z.object({
  name: z.string().min(2).max(150),
  email: z.string().email().max(150).nullable().optional(),
  phone: z.string().max(30).nullable().optional(),
  level: z.number().int().min(1).max(5),
  vehicle_type: z.string().max(50).nullable().optional(),
  license_plate: z.string().max(30).nullable().optional(),
  status: z.enum(["active", "inactive", "suspended"]).nullable().optional(),
  registered_at: z.string().refine(isDateLike, "registered_at must be a valid date or datetime").nullable().optional()
});

const allowedSort = ["name", "-name", "registered_at", "-registered_at", "created_at", "-created_at"] as const;

export function parseCourierQuery(query: Record<string, unknown>): { value?: CourierQuery; errors?: Record<string, string[]> } {
  const errors: Record<string, string[]> = {};
  const page = integerQuery(query.page, 1, 1, Number.MAX_SAFE_INTEGER, "page", errors);
  const perPage = integerQuery(query.per_page, 10, 1, 100, "per_page", errors);
  const sort = typeof query.sort === "string" && query.sort !== "" ? query.sort : "name";
  if (!allowedSort.includes(sort as CourierQuery["sort"])) {
    errors.sort = ["The sort field is invalid."];
  }

  const levels: number[] = [];
  if (typeof query.level === "string" && query.level.trim() !== "") {
    for (const part of query.level.split(",")) {
      const value = Number(part.trim());
      if (!Number.isInteger(value) || value < 1 || value > 5) {
        errors.level = ["The level query must contain only levels 1 to 5."];
        break;
      }
      levels.push(value);
    }
  }

  if (Object.keys(errors).length > 0) {
    return { errors };
  }
  return {
    value: {
      page,
      per_page: perPage,
      sort: sort as CourierQuery["sort"],
      search: typeof query.search === "string" ? query.search : undefined,
      levels
    }
  };
}

export function zodErrors(error: z.ZodError): Record<string, string[]> {
  const errors: Record<string, string[]> = {};
  for (const issue of error.issues) {
    const field = String(issue.path[0] ?? "body");
    errors[field] = [issue.message];
  }
  return errors;
}

function integerQuery(raw: unknown, fallback: number, min: number, max: number, field: string, errors: Record<string, string[]>) {
  if (raw === undefined || raw === null || raw === "") {
    return fallback;
  }
  const value = Number(raw);
  if (!Number.isInteger(value) || value < min || value > max) {
    errors[field] = [`The ${field} query is invalid.`];
    return fallback;
  }
  return value;
}

function isDateLike(value: string) {
  return !Number.isNaN(Date.parse(value));
}
