// 未登录
export class UnauthorizedError extends Error {
  code: string;
  status: number;
  constructor(message: string) {
    super(message);
    this.code = "Unauthorized";
    this.status = 401;
  }
}

// 参数错误
export class ParamError extends Error {
  code: string;
  status: number;
  constructor(message: string) {
    super(message);
    this.code = "Param Error";
    this.status = 422;
  }
}

// 提示错误
export class TipError extends Error {
  code: string;
  status: number;
  constructor(message: string) {
    super(message);
    this.code = "Tip Error";
    this.status = 427;
  }
}
