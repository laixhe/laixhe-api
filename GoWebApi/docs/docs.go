// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "登录用户",
                "parameters": [
                    {
                        "description": "请求body参数",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/refresh": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "刷新Jwt",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token令牌",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.RefreshResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "注册用户",
                "parameters": [
                    {
                        "description": "请求body参数",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterResponse"
                        }
                    }
                }
            }
        },
        "/api/user/info": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token令牌",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.InfoRequest"
                        }
                    }
                }
            }
        },
        "/api/user/list": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "每页页数(数量)",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "当前页数",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ListResponse"
                        }
                    }
                }
            }
        },
        "/api/user/update": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "请求body参数",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/user.UpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UpdateResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/user.User"
                }
            }
        },
        "auth.RefreshResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterRequest": {
            "type": "object",
            "required": [
                "age",
                "email",
                "password",
                "uname"
            ],
            "properties": {
                "age": {
                    "description": "用户年龄",
                    "type": "integer",
                    "maximum": 200,
                    "minimum": 0
                },
                "email": {
                    "description": "用户邮箱",
                    "type": "string"
                },
                "password": {
                    "description": "用户密码",
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "uname": {
                    "description": "用户名",
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 2
                }
            }
        },
        "auth.RegisterResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/user.User"
                }
            }
        },
        "user.InfoRequest": {
            "type": "object"
        },
        "user.ListResponse": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.User"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "user.UpdateRequest": {
            "type": "object",
            "required": [
                "login_at",
                "uname"
            ],
            "properties": {
                "login_at": {
                    "description": "登录时间",
                    "type": "string"
                },
                "uname": {
                    "description": "用户名",
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 2
                }
            }
        },
        "user.UpdateResponse": {
            "type": "object"
        },
        "user.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "创建时间",
                    "type": "string"
                },
                "email": {
                    "description": "用户邮箱",
                    "type": "string"
                },
                "uid": {
                    "description": "用户ID",
                    "type": "integer"
                },
                "uname": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API接口",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
