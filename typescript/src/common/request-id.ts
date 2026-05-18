import { randomUUID } from "node:crypto";
import type { FastifyReply, FastifyRequest } from "fastify";

export function requestId(request: FastifyRequest, reply: FastifyReply): string {
  const incoming = request.headers["x-request-id"];
  const value = Array.isArray(incoming) ? incoming[0] : incoming;
  const id = value && value.trim() !== "" ? value : randomUUID();
  reply.header("x-request-id", id);
  return id;
}
