import { describe, expect, it } from "vitest";
import { courierPayloadSchema, parseCourierQuery } from "./courier.schema.js";

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
