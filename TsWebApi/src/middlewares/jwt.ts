import { Elysia } from "elysia";
import jwt from "@elysiajs/jwt";

const jwtConfig = jwt({
  name: "JwtConfig",
  secret: process.env.JWT_SECRET!,
  exp: process.env.JWT_EXPIRE!,
});

const jwtAuth = (app: Elysia) =>
  app
    .use(jwtConfig)
    .derive(async ({ JwtConfig, headers: { authorization }, set }) => {
      const token = authorization?.substring(7);
      console.log("jwtAuth token", token);
      if (token) {
        const payload = await JwtConfig.verify(token);
        console.log("jwtAuth payload", payload);
        if (payload) {
          const uid = Number(payload?.uid);
          if (Number.isInteger(uid) && uid > 0) {
            return { uid };
          }
        }
      }
      set.status = 401;
      throw new Error("Unauthorized");
    });

export { jwtConfig, jwtAuth };
