import type { FastifyReply, FastifyRequest } from "fastify";
import { requestId } from "./request-id.js";

function meta(request: FastifyRequest, reply: FastifyReply) {
  return {
    request_id: requestId(request, reply),
    timestamp: new Date().toISOString()
  };
}

export function success(reply: FastifyReply, request: FastifyRequest, statusCode: number, message: string, data: unknown) {
  return reply.status(statusCode).send({ success: true, message, data, meta: meta(request, reply) });
}

export function list(
  reply: FastifyReply,
  request: FastifyRequest,
  message: string,
  data: unknown,
  pagination: { page: number; per_page: number; total: number; total_pages: number }
) {
  return reply.status(200).send({ success: true, message, data, pagination, meta: meta(request, reply) });
}

export function failure(
  reply: FastifyReply,
  request: FastifyRequest,
  statusCode: number,
  message: string,
  errors?: Record<string, string[]>
) {
  return reply.status(statusCode).send({ success: false, message, ...(errors ? { errors } : {}), meta: meta(request, reply) });
}
