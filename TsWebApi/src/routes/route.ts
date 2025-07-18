import { Elysia, ValidationError } from "elysia";

import { ParamError } from "@core/error";
import { auth } from "./auth";

export const routeV1 = new Elysia()
  .group("v1", (app) => app.use(auth))
  .onError(({ error, set }) => {
    let message = "";
    let data: any = null;
    if (set.status == 422) {
      message = "参数错误";
      try {
        if (error instanceof ValidationError) {
          let errorValue = error.validator.Errors(error.value);
          let errorValueFirst = errorValue?.First();
          let msg = errorValueFirst?.schema?.description;
          if (!msg) {
            msg = errorValueFirst?.message;
          }
          data = `${errorValueFirst?.path} : ${msg}`;
        }
        if (!data) {
          if (error instanceof Error) {
            let jsonMessage = JSON.parse(error.message);
            if (jsonMessage && jsonMessage.message && jsonMessage.on) {
              data = `${jsonMessage.on}: ${jsonMessage.property} ${jsonMessage.message}`;
            } else {
              data = error.message;
            }
          }
        }
      } catch (e) {
        data = error.toString();
      }
    } else if (error instanceof ParamError) {
      message = error.message;
    } else if (error instanceof Error) {
      message = error.message;
    } else {
      message = error.toString();
    }
    return {
      code: set.status,
      msg: message,
      data: data,
    };
  });
