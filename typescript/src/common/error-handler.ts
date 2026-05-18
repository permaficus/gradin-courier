import type { FastifyError, FastifyReply, FastifyRequest } from "fastify";
import { NotFoundError, ValidationError } from "./errors.js";
import { failure } from "./response.js";

export function errorHandler(error: FastifyError | Error, request: FastifyRequest, reply: FastifyReply) {
  if (error instanceof ValidationError) {
    return failure(reply, request, error.statusCode, "Validation failed", error.errors);
  }

  if (error instanceof NotFoundError) {
    return failure(reply, request, error.statusCode, "Courier not found");
  }

  request.log.error(error);
  return failure(reply, request, 500, "Internal server error");
}
