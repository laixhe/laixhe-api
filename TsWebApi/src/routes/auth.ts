import { Elysia } from "elysia";

import { AuthService } from "@service/auth/index";

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
  .use(AuthService.login)
  .use(AuthService.refresh);
