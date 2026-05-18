import type { FastifyReply, FastifyRequest } from "fastify";
import { list, success } from "../../common/response.js";
import type { CourierService } from "./courier.service.js";

export class CourierController {
  constructor(private readonly service: CourierService) {}

  index = async (request: FastifyRequest, reply: FastifyReply) => {
    const result = await this.service.list(request.query as Record<string, unknown>);
    return list(reply, request, "Couriers retrieved successfully", result.data, result.pagination);
  };

  store = async (request: FastifyRequest, reply: FastifyReply) => {
    const courier = await this.service.create(request.body);
    return success(reply, request, 201, "Courier created successfully", courier);
  };

  show = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    const courier = await this.service.find(request.params.id);
    return success(reply, request, 200, "Courier retrieved successfully", courier);
  };

  update = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    const courier = await this.service.update(request.params.id, request.body);
    return success(reply, request, 200, "Courier updated successfully", courier);
  };

  destroy = async (request: FastifyRequest<{ Params: { id: string } }>, reply: FastifyReply) => {
    await this.service.delete(request.params.id);
    return success(reply, request, 200, "Courier deleted successfully", null);
  };
}
