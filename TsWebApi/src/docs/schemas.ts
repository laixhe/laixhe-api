// OpenAPI components.schemas 定义
// 由 @elysiajs/swagger 插件在运行时合并进动态生成的文档 (替代原 swagger-jsdoc @openapi 注解静态生成)
// 供 src/index.ts 的 swagger({ documentation: { components: { schemas } } }) 使用

import type { OpenAPIV3 } from "openapi-types";

type Schemas = NonNullable<OpenAPIV3.Document["components"]>["schemas"];

export const schemas: Schemas = {
  Error: {
    type: "object",
    required: ["code", "message"],
    properties: {
      code: { type: "integer", description: "错误码 (与 HTTP 状态码一致)" },
      message: { type: "string", description: "错误描述" },
    },
  },
  HealthResponse: {
    type: "object",
    required: ["status", "database", "version", "started_at", "now"],
    properties: {
      status: { type: "string", description: '服务状态 (固定 "ok")' },
      database: { type: "string", description: '数据库状态 (固定 "up")' },
      version: { type: "string", description: "服务版本" },
      started_at: { type: "string", description: "服务启动时间 (服务器本地时区)" },
      now: { type: "string", description: "当前时间 (服务器本地时区)" },
    },
  },
  AuthRegisterRequest: {
    type: "object",
    required: ["nickname", "email", "password"],
    properties: {
      nickname: { type: "string", description: "昵称" },
      email: { type: "string", description: "邮箱" },
      password: { type: "string", description: "密码" },
    },
  },
  AuthLoginRequest: {
    type: "object",
    required: ["email", "password"],
    properties: {
      email: { type: "string", description: "邮箱" },
      password: { type: "string", description: "密码" },
    },
  },
  AuthTokenResponse: {
    type: "object",
    required: ["token", "user"],
    properties: {
      token: { type: "string", description: "jwt token" },
      user: { $ref: "#/components/schemas/User" },
    },
  },
  User: {
    type: "object",
    description: "用户信息",
    required: ["uid", "type_id", "account", "mobile", "email", "nickname", "avatar_url", "sex", "states", "created_at"],
    properties: {
      uid: { type: "integer", description: "用户id" },
      type_id: { type: "integer", description: "类型 1-普通用户", enum: [1] },
      account: { type: "string", description: "账号" },
      mobile: { type: "string", description: "手机号" },
      email: { type: "string", description: "邮箱" },
      nickname: { type: "string", description: "昵称" },
      avatar_url: { type: "string", description: "头像地址" },
      sex: { type: "integer", description: "性别 (0-未填写 1-男 2-女)", enum: [0, 1, 2] },
      states: { type: "integer", description: "状态 (0-禁用 1-正常)", enum: [0, 1] },
      created_at: { type: "string", description: '创建时间, 格式 "YYYY-MM-DD HH:mm:ss"' },
    },
  },
  UserUpdateRequest: {
    type: "object",
    required: ["nickname"],
    properties: {
      nickname: { type: "string", description: "昵称" },
      avatar_url: { type: "string", description: "头像地址" },
    },
  },
  UserListResponse: {
    type: "object",
    required: ["total", "page", "page_size", "list"],
    properties: {
      total: { type: "integer", description: "总数" },
      page: { type: "integer", description: "分页-当前页" },
      page_size: { type: "integer", description: "分页-每页数量" },
      list: { type: "array", description: "列表", items: { $ref: "#/components/schemas/User" } },
    },
  },
};
