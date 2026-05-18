import type { FastifyReply, FastifyRequest } from "fastify";
import { NotFoundError, ValidationError } from "../../common/errors.js";
import { failure, list, success } from "../../common/response.js";
import type { CourierService } from "./courier.service.js";

export class CourierController {
  constructor(private readonly service: CourierService) {}

  index = async (request: FastifyRequest, reply: FastifyReply) => {
    try {
      const result = await this.service.list(request.query as Record<string, unknown>);
      return list(reply, request, "Couriers retrieved successfully", result.data, result.pagination);
    } catch (error) {
      return handleError(request, reply, error, 400);
    }
  };

  store = async (request: FastifyRequest, reply: FastifyReply) => {
    try {
      const courier = await this.service.create(request.body);
      return success(reply, request, 201, "Courier created successfully", courier);
    } catch (error) {
      return handleError(request, reply, error, 422);
    }
  };

  show = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    try {
      const courier = await this.service.find(request.params.id);
      return success(reply, request, 200, "Courier retrieved successfully", courier);
    } catch (error) {
      return handleError(request, reply, error, 404);
    }
  };

  update = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    try {
      const courier = await this.service.update(request.params.id, request.body);
      return success(reply, request, 200, "Courier updated successfully", courier);
    } catch (error) {
      return handleError(request, reply, error, error instanceof NotFoundError ? 404 : 422);
    }
  };

  destroy = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    try {
      await this.service.delete(request.params.id);
      return success(reply, request, 200, "Courier deleted successfully", null);
    } catch (error) {
      return handleError(request, reply, error, 404);
    }
  };
}

function handleError(request: FastifyRequest, reply: FastifyReply, error: unknown, validationStatus: number) {
  if (error instanceof ValidationError) {
    return failure(reply, request, validationStatus, "Validation failed", error.errors);
  }
  if (error instanceof NotFoundError) {
    return failure(reply, request, 404, "Courier not found");
  }
  return failure(reply, request, 500, "Internal server error");
}
