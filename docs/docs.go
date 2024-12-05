package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "xoti$",
            "url": "https://t.me/xoticdsign",
            "email": "xoticdollarsign@outlook.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://mit-license.org/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "security": [
                    {
                        "KeyAuth": []
                    }
                ],
                "description": "Возвращает полный список цитат, хранящихся в базе данных. Полезно для получения всех доступных данных для анализа, отображения или других операций. Цитаты возвращаются в формате JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Операции с цитатами"
                ],
                "summary": "Предоставляет все цитаты",
                "operationId": "list-all",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Quote"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/random": {
            "get": {
                "security": [
                    {
                        "KeyAuth": []
                    }
                ],
                "description": "Возвращает случайную цитату из базы данных. Если цитата отсутствует в кэше, то она извлекается из базы данных, добавляется в кэш и возвращается пользователю. Позволяет отображать динамическое содержимое, не перегружая базу данных. Случайность обеспечивается генератором случайных чисел.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Операции с цитатами"
                ],
                "summary": "Предоставляет случайную цитату",
                "operationId": "random-quote",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Quote"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/{id}": {
            "get": {
                "security": [
                    {
                        "KeyAuth": []
                    }
                ],
                "description": "Возвращает цитату по её уникальному идентификатору (ID). Если цитата не найдена в кэше, происходит обращение к базе данных. Полученная цитата затем сохраняется в кэш для ускорения последующих запросов. Если запрошенного ID нет в базе данных, возвращается ошибка.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Операции с цитатами"
                ],
                "summary": "Предоставляет цитату по заданному ID",
                "operationId": "quote-id",
                "parameters": [
                    {
                        "type": "string",
                        "example": "105",
                        "description": "Позволяет указать ID цитаты",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Quote"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "responses.Quote": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "quote": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "KeyAuth": {
            "type": "apiKey",
            "name": "auf-citaty-key",
            "in": "query"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Auf Citaty API",
	Description:      "TODO",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
