import { buildApp } from "./app.js";

const port = Number(process.env.APP_PORT ?? 3000);
const app = await buildApp();

await app.listen({ host: "0.0.0.0", port });
