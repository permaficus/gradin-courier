import Fastify from "fastify";
import swagger from "@fastify/swagger";
import swaggerUi from "@fastify/swagger-ui";
import { courierRoutes } from "./modules/couriers/courier.routes.js";

export async function buildApp() {
  const app = Fastify({ logger: true });
  await app.register(swagger, {
    openapi: {
      info: { title: "Courier Technical Test API", version: "1.0.0" }
    }
  });
  await app.register(swaggerUi, { routePrefix: "/docs" });
  await app.register(courierRoutes, { prefix: "/api" });
  return app;
}
