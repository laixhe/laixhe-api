import { Elysia, ValidationError } from "elysia";

import { ParamError } from "@core/error";
import { auth } from "./auth";

export const routeV1 = new Elysia()
  .group("v1", (app) => app.use(auth))
  .onError(({ code, error, set }) => {
    let message = "";
    let data: any = null;
    if (set.status == 422) {
      message = "参数错误";
      try {
        if (error instanceof ValidationError) {
          data =
            error.validator.Errors(error.value).First().path +
            " : " +
            error.validator.Errors(error.value).First().message;
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
    console.log("route v1 onError", set.status, code, message);
    return {
      code: set.status,
      msg: message,
      data: data,
    };
  });
