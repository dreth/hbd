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
        "/add-birthday": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint adds a new birthday for the authenticated user. The request must include a valid JWT token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "birthdays"
                ],
                "summary": "Add a new birthday",
                "parameters": [
                    {
                        "description": "Add birthday",
                        "name": "birthday",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.BirthdayNameDateAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.BirthdayFull"
                        }
                    },
                    "400": {
                        "description": "Invalid request or date format",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to insert birthday",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 7
            }
        },
        "/check-birthdays": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint checks for user reminders through a POST request. The request must include a valid JWT token.",
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
        "/delete-birthday": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint deletes a birthday for the authenticated user. The request must include a valid JWT token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "birthdays"
                ],
                "summary": "Delete a birthday",
                "parameters": [
                    {
                        "description": "Delete birthday",
                        "name": "birthday",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.BirthdayNameDateModify"
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
                        "description": "Invalid request or date format",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to delete birthday",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 8
            }
        },
        "/delete-user": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint deletes a user based on their email obtained from the JWT token. The request must include a valid JWT token.",
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
                        "description": "Unauthorized",
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
        "/generate-password": {
            "get": {
                "description": "This endpoint generates a new password for the user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Generate a new password",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Password"
                        }
                    },
                    "500": {
                        "description": "Failed to generate password",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 1
            }
        },
        "/health": {
            "get": {
                "description": "This endpoint checks the readiness of the service and returns a status.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Check service readiness",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Ready"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "This endpoint logs in a user by validating their email and password. Upon successful authentication, it generates a JWT token and returns the user's details along with the filtered list of birthdays.",
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
                        "description": "Invalid email or password",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 3
            }
        },
        "/me": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint returns the authenticated user's data including Telegram bot API key, user ID, reminder time, and birthdays. The request must include a valid JWT token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Get user data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.UserData"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                }
            }
        },
        "/modify-birthday": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint modifies a birthday for the authenticated user. The request must include a valid JWT token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "birthdays"
                ],
                "summary": "Modify a birthday",
                "parameters": [
                    {
                        "description": "Modify birthday",
                        "name": "birthday",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/structs.BirthdayNameDateModify"
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
                        "description": "Invalid request or date format",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to update birthday",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    }
                },
                "x-order": 9
            }
        },
        "/modify-user": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "This endpoint modifies a user's details such as Telegram bot API key, reminder time, and more. The request must include a valid JWT token. When modifying the email or password, a new JWT token is generated and returned. Otherwise, the user's data is returned without a new token.",
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
                        "description": "User data without a new token if no email or password changes",
                        "schema": {
                            "$ref": "#/definitions/structs.UserData"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/structs.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to update user or process request",
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
                            "$ref": "#/definitions/structs.LoginSuccess"
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
        "structs.BirthdayFull": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                }
            }
        },
        "structs.BirthdayNameDateAdd": {
            "type": "object",
            "required": [
                "date",
                "name"
            ],
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                }
            }
        },
        "structs.BirthdayNameDateModify": {
            "type": "object",
            "required": [
                "date",
                "id",
                "name"
            ],
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
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
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@lotiguere.com"
                },
                "password": {
                    "type": "string",
                    "example": "9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"
                }
            }
        },
        "structs.LoginSuccess": {
            "type": "object",
            "properties": {
                "birthdays": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/structs.BirthdayFull"
                    }
                },
                "reminder_time": {
                    "type": "string",
                    "example": "15:04"
                },
                "telegram_bot_api_key": {
                    "type": "string",
                    "example": "270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"
                },
                "telegram_user_id": {
                    "type": "string",
                    "example": "123456789"
                },
                "timezone": {
                    "type": "string",
                    "example": "America/New_York"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "structs.ModifyUserRequest": {
            "type": "object",
            "required": [
                "new_reminder_time",
                "new_telegram_bot_api_key",
                "new_telegram_user_id",
                "new_timezone"
            ],
            "properties": {
                "new_email": {
                    "type": "string",
                    "example": "example2@lotiguere.com"
                },
                "new_password": {
                    "type": "string",
                    "example": "9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"
                },
                "new_reminder_time": {
                    "type": "string",
                    "example": "15:04"
                },
                "new_telegram_bot_api_key": {
                    "type": "string",
                    "example": "270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"
                },
                "new_telegram_user_id": {
                    "type": "string",
                    "example": "123456789"
                },
                "new_timezone": {
                    "type": "string",
                    "example": "America/New_York"
                }
            }
        },
        "structs.Password": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"
                }
            }
        },
        "structs.Ready": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "structs.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "reminder_time",
                "telegram_bot_api_key",
                "telegram_user_id",
                "timezone"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@lotiguere.com"
                },
                "password": {
                    "type": "string",
                    "example": "9cc76406913372c2b3a3474e8ebb8dc917bdb9c4a7c5e98c639ed20f5bcf4da1"
                },
                "reminder_time": {
                    "type": "string",
                    "example": "15:04"
                },
                "telegram_bot_api_key": {
                    "type": "string",
                    "example": "270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"
                },
                "telegram_user_id": {
                    "type": "string",
                    "example": "123456789"
                },
                "timezone": {
                    "type": "string",
                    "example": "America/New_York"
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
        },
        "structs.UserData": {
            "type": "object",
            "properties": {
                "birthdays": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/structs.BirthdayFull"
                    }
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "reminder_time": {
                    "type": "string",
                    "example": "15:04"
                },
                "telegram_bot_api_key": {
                    "type": "string",
                    "example": "270485614:AAHfiqksKZ8WmR2zSjiQ7jd8Eud81ggE3e-3"
                },
                "telegram_user_id": {
                    "type": "string",
                    "example": "123456789"
                },
                "timezone": {
                    "type": "string",
                    "example": "America/New_York"
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
