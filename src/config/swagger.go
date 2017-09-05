package config

import (
	"fmt"

	"github.com/grokify/swaggman/swagger2"
)

// http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html#api-gateway-set-up-lambda-proxy-integration-on-proxy-resource

const (
	AwsApiGatewayHostnameFormat = `%v.execute-api.%v.amazonaws.com`
	AwsApiGatewayBaseURLFormat  = `https://%v.execute-api.%v.amazonaws.com/%v/`
)

func BuildAwsApiGatewayBaseURL(restapiID string, region string, stageName string) string {
	return fmt.Sprintf(AwsApiGatewayBaseURLFormat, restapiID, region, stageName)
}

// https://{restapi_id}.execute-api.{region}.amazonaws.com/{stage_name}/
func GetSwaggerSpecAwsApiGateway(restapiID string, region string, stageName string) swagger2.Specification {
	spec := GetSwaggerSpec()
	spec.Host = fmt.Sprintf(AwsApiGatewayHostnameFormat, restapiID, region)
	spec.BasePath = fmt.Sprintf("/%v", stageName)
	return spec
}

func GetSwaggerSpec() swagger2.Specification {
	spec := swagger2.Specification{
		Swagger: "2.0",
		Info: &swagger2.Info{
			Description: "Powered by: github.com/grokify/webhookproxy",
			Version:     "v1.0.0",
		},
		Schemes: []string{"https"},
		Paths: map[string]swagger2.Path{
			"/hooks": {
				Post: &swagger2.Endpoint{
					Tags:        []string{"webhook"},
					Description: "Proxy a webhook",
					Consumes:    []string{"application/json"},
					Produces:    []string{"application/json"},
					Parameters: []swagger2.Parameter{
						{
							Description: "Format of the input message",
							In:          "query",
							Name:        "input",
							Required:    true,
							Type:        "string",
						},
						{
							Description: "Format of the output message",
							In:          "query",
							Name:        "output",
							Required:    true,
							Type:        "string",
						},
						{
							Description: "URL or UID of the output message",
							In:          "query",
							Name:        "urloruid",
							Required:    true,
							Type:        "string",
						},
						{
							Description: "Your unique token",
							In:          "query",
							Name:        "token",
							Required:    false,
							Type:        "string",
						},
					},
					Responses: map[string]swagger2.Response{
						"200": {
							Description: "Successful operation",
							Schema: swagger2.Schema{
								Ref: "#/definitions/ProxyWebhookResponse",
							},
							Headers: map[string]swagger2.Header{
								"Access-Control-Allow-Origin": {
									Type:        "string",
									Description: "URI that may access the resource"},
								"Access-Control-Allow-Methods": {
									Type:        "string",
									Description: "Method or methods allowed when accessing the resource"},
								"Access-Control-Allow-Headers": {
									Type:        "string",
									Description: "Used in response to a preflight request to indicate which HTTP headers can be used when making the request."},
								"Content-Type": {
									Type:        "string",
									Description: "MIME type of the response"},
							},
						},
					},
					XAmazonApigatewayIntegration: swagger2.XAmazonApigatewayIntegration{
						Responses: map[string]swagger2.XAmazonApigatewayIntegrationResponse{
							"default": {
								StatusCode:         "200",
								ResponseParameters: map[string]string{"method.response.header.Content-Type": "'text/html'"},
								ResponseTemplates:  map[string]string{"text/html": "<html>\n    <head>\n        <style>\n        body {\n            color: #333;\n            font-family: Sans-serif;\n            max-width: 800px;\n            margin: auto;\n        }\n        </style>\n    </head>\n    <body>\n        <h1>Welcome to your Webhook Proxy</h1>\n        <p>\n            You have successfully deployed your webhookproxy.</p>\n    </body>\n</html>"},
							}},
						PassthroughBehavior: "when_no_match",
						RequestTemplates: map[string]string{
							"application/json": "{\"statusCode\": 200}",
						},
						Type: "mock",
					},
				},
			},
		},
		Definitions: map[string]swagger2.Definition{
			"ProxyWebhookResponse": {
				Type: "object",
				Properties: map[string]swagger2.Property{
					"message": {
						Type: "string",
					},
				},
			},
		},
		XAmazonApigatewayDocumentation: swagger2.XAmazonApigatewayDocumentation{
			DocumentationParts: []swagger2.DocumentationPart{
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type: "API"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Info: &swagger2.XAmazonApigatewayDocumentationPartInfo{
							Description: "WebProxy"}},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "METHOD",
						Path:   "/hooks",
						Method: "POST"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Post a webhook"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "QUERY_PARAMETER",
						Path:   "/hooks",
						Method: "POST",
						Name:   "input"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Input style"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "QUERY_PARAMETER",
						Path:   "/hooks",
						Method: "POST",
						Name:   "output"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Output style"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "QUERY_PARAMETER",
						Path:   "/hooks",
						Method: "POST",
						Name:   "urloruid"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Output style"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "QUERY_PARAMETER",
						Path:   "/hooks",
						Method: "POST",
						Name:   "token"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Security token"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:   "REQUEST_BODY",
						Path:   "/hooks",
						Method: "POST"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Webhook object that needs to be proxied"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:       "RESPONSE",
						Method:     "*",
						StatusCode: "200"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Successful operation"},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:       "RESPONSE_HEADER",
						Method:     "OPTIONS",
						StatusCode: "200",
						Name:       "Access-Control-Allow-Headers"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Used in response to a preflight request to indicate which HTTP headers can be used when making the request."},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:       "RESPONSE_HEADER",
						Method:     "OPTIONS",
						StatusCode: "200",
						Name:       "Access-Control-Allow-Methods"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Method or methods allowed when accessing the resource."},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:       "RESPONSE_HEADER",
						Method:     "OPTIONS",
						StatusCode: "200",
						Name:       "Access-Control-Allow-Origin"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "URI that may access the resource."},
				},
				{
					Location: swagger2.XAmazonApigatewayDocumentationPartLocation{
						Type:       "RESPONSE_HEADER",
						Method:     "POST",
						StatusCode: "200",
						Name:       "Content-Type"},
					Properties: swagger2.XAmazonApigatewayDocumentationPartProperties{
						Description: "Media type of request."},
				},
			},
		},
	}

	return spec
}

/*


{
  "swagger": "2.0",
  "info": {
    "description": "Your first API with Amazon API Gateway. This is a sample API that integrates via HTTP with our demo Pet Store endpoints",
    "title": "PetStore"
  },
  "schemes": [
    "https"
  ],
  "paths": {
    "/": {
      "get": {
        "tags": [
          "pets"
        ],
        "description": "PetStore HTML web page containing API usage information",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "text/html"
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "headers": {
              "Content-Type": {
                "type": "string",
                "description": "Media type of request"
              }
            }
          }
        },
        "x-amazon-apigateway-integration": {
          "responses": {
            "default": {
              "statusCode": "200",
              "responseParameters": {
                "method.response.header.Content-Type": "'text/html'"
              },
              "responseTemplates": {
                "text/html": "<html>\n    <head>\n        <style>\n        body {\n            color: #333;\n            font-family: Sans-serif;\n            max-width: 800px;\n            margin: auto;\n        }\n        </style>\n    </head>\n    <body>\n        <h1>Welcome to your Pet Store API</h1>\n        <p>\n            You have successfully deployed your first API. You are seeing this HTML page because the <code>GET</code> method to the root resource of your API returns this content as a Mock integration.\n        </p>\n        <p>\n            The Pet Store API contains the <code>/pets</code> and <code>/pets/{petId}</code> resources. By making a <a href=\"/$context.stage/pets/\" target=\"_blank\"><code>GET</code> request</a> to <code>/pets</code> you can retrieve a list of Pets in your API. If you are looking for a specific pet, for example the pet with ID 1, you can make a <a href=\"/$context.stage/pets/1\" target=\"_blank\"><code>GET</code> request</a> to <code>/pets/1</code>.\n        </p>\n        <p>\n            You can use a REST client such as <a href=\"https://www.getpostman.com/\" target=\"_blank\">Postman</a> to test the <code>POST</code> methods in your API to create a new pet. Use the sample body below to send the <code>POST</code> request:\n        </p>\n        <pre>\n{\n    \"type\" : \"cat\",\n    \"price\" : 123.11\n}\n        </pre>\n    </body>\n</html>"
              }
            }
          },
          "passthroughBehavior": "when_no_match",
          "requestTemplates": {
            "application/json": "{\"statusCode\": 200}"
          },
          "type": "mock"
        }
      }
    },
    "/hooks": {
      "post": {
        "tags": [
          "hooks"
        ],
        "operationId": "Proxy Webhook",
        "summary": "Create a pet",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "name": "input",
            "in": "query",
            "description": "Format of the input message",
            "required": true,
            "type": "string"
          },
          {
            "name": "output",
            "in": "query",
            "description": "Format of the output message",
            "required": true,
            "type": "string"
          },
          {
            "name": "urloruid",
            "in": "query",
            "description": "URL or UID of the output webhook",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "schema": {
              "$ref": "#/definitions/HookResponse"
            },
            "headers": {
              "Access-Control-Allow-Origin": {
                "type": "string",
                "description": "URI that may access the resource"
              }
            }
          }
        },
        "x-amazon-apigateway-integration": {
          "responses": {
            "default": {
              "statusCode": "200",
              "responseParameters": {
                "method.response.header.Access-Control-Allow-Origin": "'*'"
              }
            }
          },
          "uri": "http://petstore-demo-endpoint.execute-api.com/petstore/pets",
          "passthroughBehavior": "when_no_match",
          "httpMethod": "POST",
          "type": "http"
        }
      },
      "options": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "schema": {
              "$ref": "#/definitions/Empty"
            },
            "headers": {
              "Access-Control-Allow-Origin": {
                "type": "string",
                "description": "URI that may access the resource"
              },
              "Access-Control-Allow-Methods": {
                "type": "string",
                "description": "Method or methods allowed when accessing the resource"
              },
              "Access-Control-Allow-Headers": {
                "type": "string",
                "description": "Used in response to a preflight request to indicate which HTTP headers can be used when making the request."
              }
            }
          }
        },
        "x-amazon-apigateway-integration": {
          "responses": {
            "default": {
              "statusCode": "200",
              "responseParameters": {
                "method.response.header.Access-Control-Allow-Methods": "'POST,GET,OPTIONS'",
                "method.response.header.Access-Control-Allow-Headers": "'Content-Type,X-Amz-Date,Authorization,X-Api-Key'",
                "method.response.header.Access-Control-Allow-Origin": "'*'"
              }
            }
          },
          "passthroughBehavior": "when_no_match",
          "requestTemplates": {
            "application/json": "{\"statusCode\": 200}"
          },
          "type": "mock"
        }
      }
    },
    "/pets/{petId}": {
      "get": {
        "tags": [
          "pets"
        ],
        "summary": "Info for a specific pet",
        "operationId": "GetPet",
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "description": "The id of the pet to retrieve",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "schema": {
              "$ref": "#/definitions/Pet"
            },
            "headers": {
              "Access-Control-Allow-Origin": {
                "type": "string",
                "description": "URI that may access the resource"
              }
            }
          }
        },
        "x-amazon-apigateway-integration": {
          "responses": {
            "default": {
              "statusCode": "200",
              "responseParameters": {
                "method.response.header.Access-Control-Allow-Origin": "'*'"
              }
            }
          },
          "requestParameters": {
            "integration.request.path.petId": "method.request.path.petId"
          },
          "uri": "http://petstore-demo-endpoint.execute-api.com/petstore/pets/{petId}",
          "passthroughBehavior": "when_no_match",
          "httpMethod": "GET",
          "type": "http"
        }
      },
      "options": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "description": "The id of the pet to retrieve",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "schema": {
              "$ref": "#/definitions/Empty"
            },
            "headers": {
              "Access-Control-Allow-Origin": {
                "type": "string",
                "description": "URI that may access the resource"
              },
              "Access-Control-Allow-Methods": {
                "type": "string",
                "description": "Method or methods allowed when accessing the resource"
              },
              "Access-Control-Allow-Headers": {
                "type": "string",
                "description": "Used in response to a preflight request to indicate which HTTP headers can be used when making the request."
              }
            }
          }
        },
        "x-amazon-apigateway-integration": {
          "responses": {
            "default": {
              "statusCode": "200",
              "responseParameters": {
                "method.response.header.Access-Control-Allow-Methods": "'GET,OPTIONS'",
                "method.response.header.Access-Control-Allow-Headers": "'Content-Type,X-Amz-Date,Authorization,X-Api-Key'",
                "method.response.header.Access-Control-Allow-Origin": "'*'"
              }
            }
          },
          "passthroughBehavior": "when_no_match",
          "requestTemplates": {
            "application/json": "{\"statusCode\": 200}"
          },
          "type": "mock"
        }
      }
    }
  },
  "definitions": {
    "Pets": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/Pet"
      }
    },
    "Empty": {
      "type": "object"
    },
    "NewPetResponse": {
      "type": "object",
      "properties": {
        "pet": {
          "$ref": "#/definitions/Pet"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "Pet": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "type": {
          "type": "string"
        },
        "price": {
          "type": "number"
        }
      }
    },
    "NewPet": {
      "type": "object",
      "properties": {
        "type": {
          "$ref": "#/definitions/PetType"
        },
        "price": {
          "type": "number"
        }
      }
    },
    "PetType": {
      "type": "string",
      "enum": [
        "dog",
        "cat",
        "fish",
        "bird",
        "gecko"
      ]
    }
  },
  "x-amazon-apigateway-documentation": {
    "version": "v2.1",
    "createdDate": "2016-11-17T07:03:59Z",
    "documentationParts": [
      {
        "location": {
          "type": "API"
        },
        "properties": {
          "info": {
            "description": "Your first API with Amazon API Gateway. This is a sample API that integrates via HTTP with our demo Pet Store endpoints"
          }
        }
      },
      {
        "location": {
          "type": "METHOD",
          "method": "GET"
        },
        "properties": {
          "tags": [
            "pets"
          ],
          "description": "PetStore HTML web page containing API usage information"
        }
      },
      {
        "location": {
          "type": "METHOD",
          "path": "/pets/{petId}",
          "method": "GET"
        },
        "properties": {
          "tags": [
            "pets"
          ],
          "summary": "Info for a specific pet"
        }
      },
      {
        "location": {
          "type": "METHOD",
          "path": "/pets",
          "method": "GET"
        },
        "properties": {
          "tags": [
            "pets"
          ],
          "summary": "List all pets"
        }
      },
      {
        "location": {
          "type": "METHOD",
          "path": "/pets",
          "method": "POST"
        },
        "properties": {
          "tags": [
            "pets"
          ],
          "summary": "Create a pet"
        }
      },
      {
        "location": {
          "type": "PATH_PARAMETER",
          "path": "/pets/{petId}",
          "method": "*",
          "name": "petId"
        },
        "properties": {
          "description": "The id of the pet to retrieve"
        }
      },
      {
        "location": {
          "type": "QUERY_PARAMETER",
          "path": "/pets",
          "method": "GET",
          "name": "page"
        },
        "properties": {
          "description": "Page number of results to return."
        }
      },
      {
        "location": {
          "type": "QUERY_PARAMETER",
          "path": "/pets",
          "method": "GET",
          "name": "type"
        },
        "properties": {
          "description": "The type of pet to retrieve"
        }
      },
      {
        "location": {
          "type": "REQUEST_BODY",
          "path": "/pets",
          "method": "POST"
        },
        "properties": {
          "description": "Pet object that needs to be added to the store"
        }
      },
      {
        "location": {
          "type": "RESPONSE",
          "method": "*",
          "statusCode": "200"
        },
        "properties": {
          "description": "Successful operation"
        }
      },
      {
        "location": {
          "type": "RESPONSE_HEADER",
          "method": "OPTIONS",
          "statusCode": "200",
          "name": "Access-Control-Allow-Headers"
        },
        "properties": {
          "description": "Used in response to a preflight request to indicate which HTTP headers can be used when making the request."
        }
      },
      {
        "location": {
          "type": "RESPONSE_HEADER",
          "method": "OPTIONS",
          "statusCode": "200",
          "name": "Access-Control-Allow-Methods"
        },
        "properties": {
          "description": "Method or methods allowed when accessing the resource"
        }
      },
      {
        "location": {
          "type": "RESPONSE_HEADER",
          "method": "*",
          "statusCode": "200",
          "name": "Access-Control-Allow-Origin"
        },
        "properties": {
          "description": "URI that may access the resource"
        }
      },
      {
        "location": {
          "type": "RESPONSE_HEADER",
          "method": "POST",
          "statusCode": "200",
          "name": "Content-Type"
        },
        "properties": {
          "description": "Media type of request"
        }
      }
    ]
  }
}

*/
