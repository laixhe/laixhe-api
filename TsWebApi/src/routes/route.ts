import { Elysia, ValidationError } from "elysia";

import { ParamError } from "@core/error";
import { auth } from "./auth";

export const routeV1 = new Elysia({ prefix: "api/v1" })
  .onError(({ error, set }) => {
    let code = set.status;
    let message = "";
    if (error instanceof ValidationError) {
      let path = error.valueError?.path || "";
      if (path.startsWith("/")) {
        path = path.substring(1);
      }
      message = path + ": " + (error.valueError?.message || "");
    } else if (error instanceof ParamError) {
      message = error.message;
    } else if (set.status == 422) {
      message = "参数错误";
    } else if (error instanceof Error) {
      message = error.message;
    } else {
      message = error.toString();
    }
    return {
      code: code,
      msg: message,
    };
  })
  .use(auth);
