// 参数错误
export class ParamError extends Error {
  code: string;
  status: number;
  constructor(message: string) {
    super(message);
    this.code = "PARAM_ERROR";
    this.status = 422;
  }
}

// 提示错误
export class TipError extends Error {
  code: string;
  status: number;
  constructor(message: string) {
    super(message);
    this.code = "TIP_ERROR";
    this.status = 427;
  }
}
