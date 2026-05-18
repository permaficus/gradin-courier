import type { FastifyInstance } from "fastify";
import { prisma } from "../../database/prisma.js";
import { CourierController } from "./courier.controller.js";
import { CourierRepository } from "./courier.repository.js";
import { CourierService } from "./courier.service.js";

export async function courierRoutes(app: FastifyInstance) {
  const controller = new CourierController(new CourierService(new CourierRepository(prisma)));

  app.get("/couriers", controller.index);
  app.post("/couriers", controller.store);
  app.get<{ Params: { id: string } }>("/couriers/:id", controller.show);
  app.put<{ Params: { id: string } }>("/couriers/:id", controller.update);
  app.delete<{ Params: { id: string } }>("/couriers/:id", controller.destroy);
}
