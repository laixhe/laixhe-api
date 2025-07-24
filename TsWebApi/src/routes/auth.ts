import { Elysia } from "elysia";

import { AuthController } from "@controller/auth/index";

export const auth = new Elysia({
  prefix: "auth",
  detail: {
    tags: ["Auth"],
    security: [
      {
        bearerAuth: [],
      },
    ],
  },
})
  .use(AuthController.login)
  .use(AuthController.refresh);
