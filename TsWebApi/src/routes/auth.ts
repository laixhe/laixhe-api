import { Elysia } from "elysia";

import { authLogin } from "@controller/auth/login";
import { authRefresh } from "@controller/auth/refresh";

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
  .use(authLogin)
  .use(authRefresh);
