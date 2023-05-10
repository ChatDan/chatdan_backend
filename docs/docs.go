// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://swagger.io/terms/",
        "contact": {
            "name": "JingYiJun",
            "url": "https://www.jingyijun.xyz",
            "email": "jingyijun3104@outlook.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/messageBox": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MessageBox Module"
                ],
                "summary": "创建提问箱",
                "parameters": [
                    {
                        "description": "box",
                        "name": "box",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.BoxCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.BoxCommonResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/messageBox/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MessageBox Module"
                ],
                "summary": "获取提问箱信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.BoxGetResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MessageBox Module"
                ],
                "summary": "修改提问箱信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "box",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.BoxModifyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.BoxCommonResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MessageBox Module"
                ],
                "summary": "删除提问箱",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/models.Map"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/messageBoxes": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MessageBox Module"
                ],
                "summary": "查询提问箱",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "integer",
                        "name": "owner",
                        "in": "query"
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "title",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.BoxListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/post": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post Module"
                ],
                "summary": "创建帖子",
                "parameters": [
                    {
                        "description": "post",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.PostCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.PostCommonResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/post/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post Module"
                ],
                "summary": "获取帖子信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.PostGetResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Only the owner of library can delete it",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post Module"
                ],
                "summary": "删除帖子",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "post",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.PostModifyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.PostCommonResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Only the owner of library can delete it",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post Module"
                ],
                "summary": "删除帖子",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Only the owner of library can delete it.",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.PostDeleteResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/posts": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post Module"
                ],
                "summary": "查询所有帖子",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "message_box_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.PostListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Module"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "The two fields are required, you can also add other fields(e.g. email).",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "用户名或密码错误",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Module"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "The two fields are required, you can also add other fields(e.g. email).",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/user/reset": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Module"
                ],
                "summary": "Reset Password",
                "parameters": [
                    {
                        "description": "json",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apis.ResetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid JWT Token",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/wall": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wall Module"
                ],
                "summary": "获取今日表白墙",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/apis.WallListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apis.BoxCommonResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "owner_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "apis.BoxCreateRequest": {
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "title": {
                    "type": "string"
                }
            }
        },
        "apis.BoxGetResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "owner_id": {
                    "type": "integer"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "apis.BoxListResponse": {
            "type": "object",
            "properties": {
                "messageBoxes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/apis.BoxCommonResponse"
                    }
                }
            }
        },
        "apis.BoxModifyRequest": {
            "type": "object",
            "properties": {
                "title": {
                    "type": "string"
                }
            }
        },
        "apis.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "minLength": 2
                }
            }
        },
        "apis.PostCommonResponse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "poster_id": {
                    "type": "integer"
                },
                "visibility": {
                    "type": "string"
                }
            }
        },
        "apis.PostCreateRequest": {
            "type": "object",
            "required": [
                "content",
                "message_box_id"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "message_box_id": {
                    "type": "integer"
                },
                "visibility": {
                    "type": "string",
                    "default": "public",
                    "enum": [
                        "public",
                        "private"
                    ]
                }
            }
        },
        "apis.PostDeleteResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "apis.PostGetResponse": {
            "type": "object",
            "properties": {
                "channels": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "poster_id": {
                    "type": "integer"
                },
                "visibility": {
                    "type": "string"
                }
            }
        },
        "apis.PostListResponse": {
            "type": "object",
            "properties": {
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/apis.PostCommonResponse"
                    }
                }
            }
        },
        "apis.PostModifyRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string",
                    "maxLength": 2000,
                    "minLength": 1
                },
                "visibility": {
                    "type": "string",
                    "enum": [
                        "public",
                        "private"
                    ]
                }
            }
        },
        "apis.ResetRequest": {
            "type": "object",
            "properties": {
                "new_password": {
                    "type": "string",
                    "minLength": 8
                },
                "old_password": {
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "apis.UserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "apis.WallListResponse": {
            "type": "object",
            "properties": {
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/apis.PostCommonResponse"
                    }
                }
            }
        },
        "models.Map": {
            "type": "object",
            "additionalProperties": {}
        },
        "utils.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "error_msg": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "用户模块",
            "name": "User Module"
        },
        {
            "description": "提问箱模块",
            "name": "MessageBox Module"
        },
        {
            "description": "帖子模块",
            "name": "Post Module"
        },
        {
            "description": "表白墙模块",
            "name": "Wall Module"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "localhost:8000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "ChatDan Backend",
	Description:      "ChatDan, a message box and 'biaobai' platform for Fudaners.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
