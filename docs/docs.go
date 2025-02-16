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
        "/auth/login": {
            "post": {
                "description": "Учетные записи администратора создаются через cli утилиту",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Вход в учетную запись администратора",
                "parameters": [
                    {
                        "description": "Данные для входа",
                        "name": "input",
                        "in": "body",
                        "required": true,
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
                    },
                    "400": {
                        "description": "Неверный формат данных",
                        "schema": {
                            "$ref": "#/definitions/httpx.Response"
                        }
                    },
                    "401": {
                        "description": "Неверные данные для входа",
                        "schema": {
                            "$ref": "#/definitions/httpx.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/httpx.Response"
                        }
                    }
                }
            }
        },
        "/content/generate": {
            "get": {
                "description": "Генерирует контент для телеграм поста на заданную тему с помощью AI",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "content"
                ],
                "summary": "Генерация контента для поста",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Тема контента",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/content.GenerateContentResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "$ref": "#/definitions/httpx.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/httpx.Response"
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
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "content.GenerateContentResponse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/httpx.Status"
                }
            }
        },
        "httpx.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/httpx.Status"
                }
            }
        },
        "httpx.Status": {
            "type": "string",
            "enum": [
                "success",
                "error"
            ],
            "x-enum-varnames": [
                "StatusSuccess",
                "StatusError"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "FitFlow API",
	Description:      "Описание API для сервиса FitFlow",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
