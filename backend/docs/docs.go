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
        "/check-birthdays": {
            "post": {
                "description": "This endpoint checks for user reminders through a POST request.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reminders"
                ],
                "summary": "Check user reminders",
                "parameters": [
                    {
                        "description": "Check reminders",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Success"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Error querying users",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 6
            }
        },
        "/delete-user": {
            "delete": {
                "description": "This endpoint deletes a user based on their email and encryption key.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "description": "Delete user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Success"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "401": {
                        "description": "Invalid encryption key or email",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to delete user",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 5
            }
        },
        "/generate-encryption-key": {
            "get": {
                "description": "This endpoint generates a new encryption key for the user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Generate a new encryption key",
                "responses": {
                    "200": {
                        "description": "encryption_key",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to generate encryption key",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 1
            }
        },
        "/login": {
            "post": {
                "description": "This endpoint logs in a user by validating their email and encryption key.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "Login user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.LoginSuccess"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "401": {
                        "description": "Invalid encryption key or email",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 3
            }
        },
        "/modify-user": {
            "put": {
                "description": "This endpoint modifies a user's details such as Telegram bot API key, reminder time, and more.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Modify a user's details",
                "parameters": [
                    {
                        "description": "Modify user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.ModifyUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Success"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "401": {
                        "description": "Invalid encryption key",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to update user",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 4
            }
        },
        "/register": {
            "post": {
                "description": "This endpoint registers a new user with their email, Telegram bot API key, and other details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Register user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Success"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "409": {
                        "description": "Email or Telegram bot API key already registered",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to create user",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 2
            }
        }
    },
    "definitions": {
        "structs.Birthday": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "structs.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "structs.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "encryption_key"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "encryption_key": {
                    "type": "string"
                }
            }
        },
        "structs.LoginSuccess": {
            "type": "object",
            "properties": {
                "birthdays": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/structs.Birthday"
                    }
                },
                "reminder_time": {
                    "type": "string"
                },
                "telegram_bot_api_key": {
                    "type": "string"
                },
                "telegram_user_id": {
                    "type": "string"
                },
                "timezone": {
                    "type": "string"
                }
            }
        },
        "structs.ModifyUserRequest": {
            "type": "object",
            "required": [
                "encryption_key",
                "reminder_time",
                "telegram_bot_api_key",
                "telegram_user_id",
                "timezone"
            ],
            "properties": {
                "birthdays": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/structs.Birthday"
                    }
                },
                "email": {
                    "type": "string"
                },
                "encryption_key": {
                    "type": "string"
                },
                "reminder_time": {
                    "type": "string"
                },
                "telegram_bot_api_key": {
                    "type": "string"
                },
                "telegram_user_id": {
                    "type": "string"
                },
                "timezone": {
                    "type": "string"
                }
            }
        },
        "structs.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "encryption_key",
                "reminder_time",
                "telegram_bot_api_key",
                "telegram_user_id",
                "timezone"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "encryption_key": {
                    "type": "string"
                },
                "reminder_time": {
                    "type": "string"
                },
                "telegram_bot_api_key": {
                    "type": "string"
                },
                "telegram_user_id": {
                    "type": "string"
                },
                "timezone": {
                    "type": "string"
                }
            }
        },
        "structs.Success": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean"
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
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}