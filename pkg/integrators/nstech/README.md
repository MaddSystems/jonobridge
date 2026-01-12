# Nstech 

📘 INTEGRATIONS HUB - NSTECH


### 🧩 Transaction Identification

Whenever a command is sent to the Hub via the API, the request response will include a unique UUID for the transaction.

This UUID is essential for tracking the command, as it serves as a unique identifier key. It should be used later to correlate the command delivery status events returned by the target technology.

🔁 Why is this important?
Tracking the command status (success, error, pending, etc.) depends on this UUID.

📘 For more details on how to handle these status events, consult the event documentation.


## 🚀 Sending Commands via the Hub

### 🧩 What are commands?

Commands are messages sent to devices of integrated technologies (such as trackers, sensors, etc.). The origin of these messages is the end application, which sends the command to the Hub — responsible for translating it and routing it correctly to the target technology.

This flow ensures standardization, traceability, and centralized integration across different platforms.

### 📤 How to send a command via the Hub?

Commands are sent via a REST request to a specific Hub endpoint, depending on the environment (staging or production).

### 🌐 Available endpoints

- Staging environment (STG):
  https://stg.nstech.com.br/zeus/api/orchestrator
- Production environment (PROD):
  https://zeus.nstech.com.br/api/orchestrator

### 📤 Request format

Commands must follow the JSON Schema defined by the Hub. This schema ensures the message structure is correct for each command type and technology. You can view, test, and understand all expected fields in our API Reference: [API Reference](https://hub-nstech.readme.io/reference/post_v1-commands)

### Example OpenAPI JSON

```json
{
  "openapi": "3.0.1",
  "info": {
    "title": "Hub - Comandos",
    "version": "v1"
  },
  "servers": [
    {
      "url": "https://stg.nstech.com.br/zeus/api/orchestrator",
      "description": "Staging Environment"
    },
    {
      "url": "<https://zeus.nstech.com.br/api/orchestrator>",
      "description": "Production Environment"
    }
  ],
  "paths": {
    "/v1/commands": {
      "post": {
        "tags": [
          "Commands"
        ],
        "summary": " SendCommand",
        "parameters": [
          {
            "name": "api-version",
            "in": "query",
            "schema": {
              "type": "string",
 		  "default": "v1"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json-patch+json": {
              "schema": {
                "$ref": "#/components/schemas/SendTechnologyCommandCommand"
              }
            },
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SendTechnologyCommandCommand"
              }
            },
            "text/json": {
              "schema": {
                "$ref": "#/components/schemas/SendTechnologyCommandCommand"
              }
            },
            "application/*+json": {
              "schema": {
                "$ref": "#/components/schemas/SendTechnologyCommandCommand"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Success",
            "content": {
              "text/plain": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "format": "uuid"
                  }
                }
              },
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "format": "uuid"
                  }
                }
              },
              "text/json": {
                "schema": {
                  "type": "array",
              }
            }
          },
          "400": {
            "description": "Bad Request",
            "content": {
              "text/plain": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "text/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              }
            }
          },
          "422": {
            "description": "Client Error",
            "content": {
              "text/plain": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "text/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              }
            }
          },
          "500": {
            "description": "Server Error",
            "content": {
              "text/plain": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              },
              "text/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProblemDetails"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "CommandType": {
        "enum": [
          "SinglePosition",
          "LockVehicle",
          "UnlockVehicle",
          "ActivateTrackingMode",
          "ActivateInteractiveMode",
          "Deactivate",
          "SendMessage",
          "SendBeep",
          "SealTrunk",
          "UnlockTrunk",
          "UnsealCabin",
          "SealCabin",
          "SealTrailer",
          "LockDoor",
          "LockTrunk",
          "UnsealTrailer",
          "SingleTemperature",
          "SealEngine",
          "UnsealEngine",
          "AuthorizeArrival",
          "AssociateTargets",
          "TransmitTargets",
          "ExcludeTargets",
          "AssociateOperation",
          "TransmitOperation",
          "ProhibitTrailerUnhitching",
          "ProhibitDoors",
          "TransitionsSendConfiguration",
          "TransitionsRequestShippingStatus",
          "TransitionsReleaseFPs",
          "TransitionsinhibitEraseFPs",
          "ImportTargets",
          "ListTargets",
          "AuthorizeOpeningTrunk",
          "AuthorizeUnsealingTrunk",
          "AuthorizeUncoupleTrailer",
          "DeclineAuthorization",
          "AuthorizeUnsealingCabin",
          "AuthorizeConnectionWithCentral",
          "AuthorizeUnsealEngine",
          "EmbarkingDriver",
          "DisembarkingDriver",
          "AuthorizeDriverDoor",
          "AuthorizePassengerDoor",
          "UnauthorizeDriverDoor",
          "UnauthorizePassengerDoor",
          "UnlockTrunkEnableLockButton",
          "LockTrunkDisableLockButton",
          "ChangePF6Minutes",
          "ChangePF10Minutes",
          "ChangePF15Minutes",
          "ChangePF20Minutes",
          "ChangePF30Minutes",
          "ChangePF60Minutes",
          "TurnOnSiren",
          "TurnOffSiren",
          "TurnOnWarning",
          "TurnOffWarning",
          "TurnOnHazardLight",
          "TurnOffHazardLight",
          "EmbarkingGroupMacro",
          "ClearGroupMacro",
          "ActivateGroupMacro",
          "ClearOperationProfile",
          "ClearChekpointProject",
          "LinkTruckTrailer",
          "UnlinkTruckTrailer",
          "EmbarkingCheckpointProject",
          "EmbarkingOperationProfile",
          "ManeuveringRadiusConfig",
          "DisableFifthWheelButton",
          "EnableFifthWheelButton",
          "LockFifthWheel",
          "UnlockFifthWheel",
          "SendMacro",
          "EmbarkingSequencing",
          "EmbarkingLayoutTD50",
          "EmbarkingLayoutEmbarkedAction",
          "DisembarkingLayoutActionAVD",
          "EmbarkingLayoutActionAVD",
          "DisembarkingLayoutAreaGroupAVD",
          "ResetAlarm",
          "SafeMode",
          "EmbarkingLayoutPointGroup",
          "DesembarkingLayoutPointGroup",
          "EmbarkingDailyPoint",
          "SensorTypes",
          "ConfigurationType",
          "ActuatorsType",
          "MacroType"
        ],
        "type": "string"
      },
      "ProblemDetails": {
        "type": "object",
        "properties": {
          "type": {
            "type": "string",
            "nullable": true
          },
          "title": {
            "type": "string",
            "nullable": true
          },
          "status": {
            "type": "integer",
            "format": "int32",
            "nullable": true
          },
          "detail": {
            "type": "string",
            "nullable": true
          },
          "instance": {
            "type": "string",
            "nullable": true
          },
          "extensions": {
            "type": "object",
            "additionalProperties": { },
            "nullable": true,
            "readOnly": true
          }
        },
        "additionalProperties": false
      },
      "SendTechnologyCommandCommand": {
        "required": [
          "account_id",
          "command_type",
          "technology_id"
        ],
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "technology_id": {
            "minLength": 1,
            "type": "string",
            "format": "uuid"
          },
          "command_type": {
            "$ref": "#/components/schemas/CommandType"
          },
          "device_id": {
            "type": "string",
            "nullable": true
          },
          "json_content_type": {
            "type": "string",
            "nullable": true
          },
          "json_content": {
            "type": "string",
            "nullable": true
          },
          "account_id": {
            "minLength": 1,
            "type": "string",
            "format": "uuid"
          },
          "extra_data": {
            "type": "object",
            "additionalProperties": { },
            "nullable": true
          },
          "event_type": {
            "$ref": "#/components/schemas/TechnologyEventType"
          }
        },
        "additionalProperties": false
      },
      "TechnologyEventType": {
        "enum": [
          "TechnologyException",
          "LogisticEvent",
          "ConfirmTransmitOperation",
          "ConfirmTransmitTargets",
          "ConfirmAssociationTargets",
          "ConfirmAssociationOperation",
          "ConfirmStatusTeleCommand",
          "UpdateDeviceInformation",
          "TargetProfile",
          "ImportationOfPredefinedMacros",
          "ConfirmReadMsgFree",
          "Temperature",
          "ArrivalTargetArea",
          "ExitTargetArea",
          "TargetArrivalCommand",
          "TargetExitCommand",
          "StartOfService",
          "EndOfService",
          "RightDoorOpening",
          "RightDoorClosing",
          "LeftDoorClosing",
          "LeftDoorOpening",
          "IgnitionOff",
          "IgnitionOn",
          "ReceivedFormattedMessage",
          "OperationReceived",
          "ReceivedTargetProfile",
          "CommandStatusReceived",
          "Authorization",
          "ReceivedPredefinedMessage",
          "InvalidXml",
          "PositionReceived",
          "ConfirmAssociationTargetsFail",
          "MessageReceived",
          "ReceivedVehicleAccessory",
          "ReceivedEmbarkedMacro",
          "ReceivedGroupMacro",
          "ReceivedItemMacro",
          "ReceivedOperationalProfile",
          "ReceivedVehicle",
          "ReceivedTrailerTruck",
          "ReceivedRedundantVehicle",
          "ReceivedCheckpoint",
          "StatusVehicle",
          "Position",
          "ReceivedProfileConfiguration"
        ],
        "type": "string"
      }
    },
    "securitySchemes": {
      "Bearer": {
        "type": "http",
        "description": "Enter JWT Bearer token **_only_**",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  },
  "security": [
    {
      "Bearer": [ ]
    }
  ]
}
```

## Overview

## General information

During the development of the INTEGRATIONS HUB we identified a significant challenge: many tracking technologies, carriers and risk managers could not be quickly benefited by point-to-point integrations.

The large diversity of these technologies — each with its own particularities and specific characteristics — made individual integrations impractical and non-competitive. Some technologies also lack proper structure or documentation, which further complicates the integration work.

As we expanded our scope to Latin America, which is a core part of our portfolio, this complexity increased. It became essential to look for solutions that enable efficient and comprehensive integration while taking into account the region's specific challenges.

As a response, we developed INTEGRA NSTECH — an innovative tool that enables reception of logistics information from satellite and tracking technologies. The platform efficiently forwards that information to interested parties, including risk managers, fleet managers, carriers and other companies operating in the same niche.

INTEGRA NSTECH provides a fast and broad solution to tackle the challenges caused by the diversity of technologies in Latin America. The platform simplifies and optimizes integration processes and offers an effective, competitive alternative for organizations seeking efficient, region-adapted logistics management.

INTEGRA NSTECH uses a unified communication language

Tracking technologies that send data according to the established standard become automatically compatible with all companies that use this tool. This removes the need for individual integrations and brings a significant advantage to both tracking technology vendors and the risk managers and carriers already integrated with the Integrations HUB.

This simplified approach not only speeds up the onboarding of new technologies but also allows risk managers and carriers connected to the HUB to access data from those new technologies with minimal additional effort. INTEGRA NSTECH aims to optimize interoperability in the logistics ecosystem, promoting efficient integration that benefits all participants.

In short, INTEGRA NSTECH does not compete with the Integrations HUB but works as a complementary solution. While the HUB actively "pulls" information from technologies after an extensive integration process — involving analysis, development, testing and technical validation — INTEGRA provides a passive communication channel. This enables tracking technologies to proactively send data, simplifying the process and eliminating the need for extensive intervention.

INTEGRA NSTECH therefore stands out as an efficient alternative, offering a more direct and agile way to obtain logistics data and complementing the Integrations HUB. Together these tools improve efficiency and flexibility in the logistics ecosystem, benefiting companies, risk managers and carriers within the NSTECH environment.

During the development of NSTECH's Integrations HUB, we encountered a significant challenge: the large diversity of technologies used by carriers and risk managers across Latin America.

Each technology has its own unique particularities, and many of them lack solid technical foundations or adequate documentation. This makes individual integrations slow, costly, and less competitive.

As we expanded our scope to the entire Latin American region, we realized the traditional point-to-point integration model would be unfeasible. We needed a more efficient and scalable approach capable of handling this growing complexity.

## Proposed Solution

The Integrations HUB was designed as a centralized infrastructure that enables unified communication between different systems. Through the HUB, information can be retrieved from integrated technologies after a technical process that includes:

- Study and certification (homologation) of the technology,
- Development of the specific integration,
- Joint testing and validation.

This model provides a robust and secure integration, making it ideal for technologies with solid technical foundations and for companies seeking a stable, long-term solution.

## Role of the HUB

- Act as a central communication point for companies within the logistics ecosystem;
- Actively retrieve information from the integrated technologies;
- Serve as the foundation for robust and stable integrations.

The HUB is best suited for situations where the involved parties have the technical availability to carry out a complete integration, even if the implementation takes more time.

## 🔐 HUB Authentication

This section gathers the initial steps and essential guidance for integrating your application with the Hub APIs securely.

### Hub access credentials

To ensure security and access control, all requests to the Hub APIs require authentication. First, you must obtain your access credentials, which act as a unique "identification key" for your company.

The credentials consist of:

- client_id: your unique identifier in the Hub
- client_secret: your secret key (keep it safe!)

These values are used to generate an authorization token, which is required to authenticate any API call.

### How to request your credentials

Send an email to erivelton.anjos@nstech.com.br with the following information:

- Name of the product or application that will use the APIs.
- Your company name.
- Your company's CNPJ.

### Remember

Credentials do not directly access the APIs — they are used to generate a token. That token must be sent in the Authorization header of each request in the following format:

{ "headers": { "Authorization": "Bearer {your_token}" } }

### Obtaining the token

For the complete step-by-step instructions on how to generate the token using your credentials, click here (link to the authentication section).

### Example: Obtain a token (Python)

Below is a minimal Python example using `requests` to call the token endpoint. You may need to include form data (for example `client_id`, `client_secret` and `grant_type`) depending on your authentication setup.

```python
import requests

url = "https://dev.nstech.com.br/auth/realms/zeus-stg/protocol/openid-connect/token"

headers = {
    "accept": "application/json",
    "content-type": "application/x-www-form-urlencoded"
}

response = requests.post(url, headers=headers)

print(response.text)
```

## Authentication — Obtaining an authorization token for the APIs

Once you have your client credentials, the next step is to obtain an authorization token. Make a POST request to the Hub nstech authentication service at the endpoint that matches your environment.

- Staging (STG):

  https://dev.nstech.com.br/auth/realms/zeus-stg/protocol/openid-connect/token

- Production (PROD):

  https://auth.nstech.com.br/realms/zeus/protocol/openid-connect/token

Include the following form fields in the POST request (application/x-www-form-urlencoded):

- client_id — the client ID provided when you obtained credentials.
- client_secret — the client secret provided when you obtained credentials.
- grant_type — must be set to "client_credentials" to use the client credentials flow.

After a successful request you will receive a response containing the authorization token (access_token). Store and manage this token securely. Use it in subsequent Hub API requests by adding the Authorization header:

```bash
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json
```

Or a minimal curl example to obtain a token (staging):

```bash
curl -X POST \
  https://dev.nstech.com.br/auth/realms/zeus-stg/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET&grant_type=client_credentials"
```

Token expiration

Authorization tokens have a limited lifetime. When the token expires you must request a new one using the same process.

Recommendation

- Reuse the same access token for all requests while it is valid. This reduces resource usage and improves efficiency and security.
- Implement logic to track the token expiration and proactively renew the token when it is expired or close to expiring.

## Example: Sniper integration (staging)

You provided a Sniper integration sample for the staging environment. Below is a safe, non-secret example of how to use those values. Do NOT commit client secrets to the repository — store them in environment variables or a secrets manager.

Provided values (example):
- TechnologyId: `52466691-482f-48a0-adfc-a68e776eb966`
- Account_id: `52a4b1da-8e17-49c5-b490-d98ff1b390e0`
- DeviceID: the device ID must be supplied by you; submissions without it will be discarded.

Authentication endpoint (staging):
https://dev.nstech.com.br/auth/realms/zeus-stg/protocol/openid-connect/token

Environment variables (recommended):
```bash
export NSTECH_CLIENT_ID="52466691-482f-48a0-adfc-a68e776eb966"
export NSTECH_CLIENT_SECRET="<your-client-secret>"
export NSTECH_TOKEN_URL="https://dev.nstech.com.br/auth/realms/zeus-stg/protocol/openid-connect/token"
```

Minimal curl to obtain a token (reads secret from env):
```bash
curl -s -X POST "$NSTECH_TOKEN_URL" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=$NSTECH_CLIENT_ID&client_secret=$NSTECH_CLIENT_SECRET&grant_type=client_credentials"
```

Minimal Python example (requests) to obtain token:
```python
import os
import requests

token_url = os.environ.get('NSTECH_TOKEN_URL')
data = {
    'client_id': os.environ.get('NSTECH_CLIENT_ID'),
    'client_secret': os.environ.get('NSTECH_CLIENT_SECRET'),
    'grant_type': 'client_credentials'
}
resp = requests.post(token_url, data=data)
print(resp.json())
```

Example: POST a command to the Hub (use the access_token returned by the previous request):
```bash
ACCESS_TOKEN="<access_token>"
curl -X POST "https://stg.nstech.com.br/zeus/api/orchestrator/v1/commands?api-version=v1" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "technology_id": "52466691-482f-48a0-adfc-a68e776eb966",
    "account_id": "52a4b1da-8e17-49c5-b490-d98ff1b390e0",
    "device_id": "<your_device_id>",
    "command_type": "SinglePosition"
  }'
```

Security note:
- If the client secret shown to you was accidentally committed to any public place, rotate it immediately.
- Prefer using a secrets manager or environment variables in CI/CD pipelines.



## Integration interfaces

The API (Application Programming Interface) is a programming standard composed of sets of instructions intended to simplify integration between different software platforms — for example, the Integrations HUB, INTEGRA NSTECH, and other tools you may use.

This standardized interface enables different systems to communicate and share data efficiently, promoting interoperability between platforms.

We use a REST architecture for our APIs, served over HTTPS and aligned with standard RESTful best practices.

Responses and initial validations follow standard HTTP status code semantics:

- 1xx: Informational — Request received, processing continues;
- 2xx: Success — The request was received, understood and accepted;
- 3xx: Redirection — Further action is required to complete the request;
- 4xx: Client Error — The request contains malformed syntax or cannot be fulfilled;
- 5xx: Server Error — The server failed to fulfill an apparently valid request.

After these initial validations, command returns, messages and handled errors are published for consumption via messaging (RabbitMQ). A dedicated queue is created exclusively for each account integrated with the platform, ensuring efficient and targeted communication.


## Technologies
## Operation flow

Operation flow — INTEGRA NSTECH

One of INTEGRA NSTECH's main characteristics is the passive interaction model: the technology initiates the process.

The first step, performed asynchronously, is obtaining access credentials (see the step-by-step instructions for obtaining credentials). This credential, which we call the Token, is used for all operations in the platform.

The Token is reusable

You don't need to request a new token for every request — the same token used in one request can be reused until it expires. Request a new token only when the current one has expired.

With a valid access credential ready, the next step is to prepare the package with the information and send it through the API. Each type of information has a specific endpoint for submission, for example:

- CommandStatus (API playground)
- Events (API playground)
- Messages (API playground)
- Positions (API playground)
- Temperatures (API playground)

After submission, the API performs initial validation on required fields and content and quickly returns confirmation indicating whether the information was processed.

APIs are type-specific

Each API is prepared to receive a specific type of information. Sending data to the wrong endpoint/queue will generate an error reported in the request response.

Validation error responses will be published to an exclusive queue and follow the schema below:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "type": {
      "type": "string",
      "nullable": true,
      "description": "A URI reference [RFC3986] that identifies the problem type. When dereferenced it should provide human-readable documentation for the problem type. When not present, the value is assumed to be 'about:blank'."
    },
    "title": {
      "type": "string",
      "nullable": true,
      "description": "A short, human-readable summary of the problem type. It SHOULD NOT change between occurrences except for localization."
    },
    "status": {
      "type": "integer",
      "nullable": true,
      "description": "The HTTP status code generated by the origin server for this occurrence of the problem."
    },
    "detail": {
      "type": "string",
      "nullable": true,
      "description": "A human-readable explanation specific to this occurrence of the problem."
    },
    "instance": {
      "type": "string",
      "nullable": true,
      "description": "A URI reference that identifies the specific occurrence of the problem."
    },
    "extensions": {
      "type": "object",
      "description": "Problem type definitions MAY extend the problem details object with additional members.",
      "additionalProperties": {
        "type": "string"
      }
    }
  }
}
```

This completes the submission process. Now, regarding information received from technologies: incoming data may be one of four types:

- Commands
- Messages
- Errors
- Responses

You can access the Hub's command submission API Reference to test and obtain packages in your receiving queue.

The consumer side uses queues

Messaging is performed using RabbitMQ. Retrieving these messages is different from sending — there is no API for consumption. An exclusive queue is created for each company to consume the information.

Each company should define a polling frequency to retrieve data from its queue.

This queue is specific for receiving information; any data sent to the platform must be sent via the API.

Below is a table with the identification values for technologies in the Staging and Production environments.

| Technology | Id (Staging) | Id (Production) |
|---|---|---|
| Omnilink | d261e78f-2d2f-427f-ad02-1f9af2f46ff0 | 27171e95-3de8-4ff2-9212-7281e8f9d725 |
| Trucks Control (OnixSat) | a1aaf8f8-905d-457d-acd9-e919a598caf2 | d54b9bec-82f8-40f0-8d11-7bfad5c565e1 |
| Onix | A1AAF8F8-905D-457D-ACD9-E919A598CAF2 | d54b9bec-82f8-40f0-8d11-7bfad5c565e1 |
| Sascar | C511A3C8-6E58-4D0B-9B05-7C1E8C82B67D | ff90a32c-de24-4c9b-b14e-c395766113ea |
| Autotrac | 685AD186-5E8F-4251-B112-A1B56A65D6BF | bc62c83a-615a-486a-8fbb-25c6369f7042 |
| Ravex | 2533AC71-6D29-4FB6-A04D-EB982754C954 | 82c30ce4-dbda-4aa1-b9f5-a0de8c996737 |
| Positron | cc624d46-0c6c-430d-977d-9d9e4315a953 | — |
| Bermann | — | 8159d490-7037-40dc-87a8-7f1f29cf01ee |

> Note: an em dash (—) indicates that no ID was provided for that environment.

## Integrated Technologies

| Name | Products | Not Integrated Products |
|---|---|---|
| Autotrac | — | — |
| Omnilink | OmniDual (RI 4454)<br>OmniTurbo (RI 4484)<br>OmniFrota<br>RI1454<br>RI1464<br>RI1484 | Omniweb (Linker) |
| Sascar | — | — |
| Trucks Control (OnixSat) | — | spy |
| Onix | — | — |
| Ravex | — | — |

## ✅ EventType

The `EventType` field represents the type of event being processed or received by the system. Each event type has a unique numeric identifier, which is used to categorize and interpret the event correctly during integration.

### EventType mapping (JSON)

The following JSON shows the numeric identifiers for each `EventType` used by the system:

```json
{
  "EventType": {
    "technology_exception": 1,
    "logistic_event": 2,
    "confirm_transmit_operation": 3,
    "confirm_transmit_targets": 4,
    "confirm_association_targets": 5,
    "confirm_association_operation": 6,
    "confirm_status_tele_command": 7,
    "update_device_information": 8,
    "target_profile": 9,
    "importation_of_predefined_macros": 10,
    "confirm_read_msg_free": 11,
    "temperature": 12,
    "arrival_target_area": 13,
    "exit_target_area": 14,
    "target_arrival_command": 15,
    "target_exit_command": 16,
    "start_of_service": 17,
    "end_of_service": 18,
    "right_door_opening": 19,
    "right_door_closing": 20,
    "left_door_closing": 21,
    "left_door_opening": 22,
    "ignition_off": 23,
    "ignition_on": 24,
    "received_formatted_message": 25,
    "operation_received": 26,
    "received_target_profile": 27,
    "command_status_received": 28,
    "authorization": 29,
    "received_predefined_message": 30,
    "invalid_xml": 31,
    "position_received": 32,
    "confirm_association_targets_fail": 33,
    "message_received": 34,
    "received_vehicle_accessory": 35,
    "received_embarked_macro": 36,
    "received_group_macro": 37,
    "received_item_macro": 38,
    "received_operational_profile": 39,
    "received_vehicle": 40,
    "received_trailer_truck": 41,
    "received_redundant_vehicle": 42,
    "received_checkpoint": 43,
    "status_vehicle": 44,
    "position": 45,
    "received_profile_configuration": 46
  }
}
```

## Schemas

Matrix of Event Schemas by Technology

The matrix below lists the event schemas available in the Hub for each technology. An "x" indicates that the schema is available for that technology.

| Schema | Omnilink | Onix | Sascar | Autotrac | Ravex |
|---|:---:|:---:|:---:|:---:|:---:|
| AuthorizationReceived | x |  |  |  |  |
| CommandStatusReceived | x |  |  |  |  |
| ConfirmAssociationOperation | x |  |  |  |  |
| ConfirmAssociationTargets | x |  |  |  |  |
| ConfirmStatusTeleCommand | x | x | x | x | x |
| ConfirmTransmitOperation | x |  |  |  |  |
| ConfirmTransmitTargets | x |  |  |  |  |
| Event (SchemaBase) | x | x | x | x | x |
| EventMessage | x | x | x | x | x |
| FormattedMessage | x |  |  |  |  |
| FormattedmessagePredefined | x |  |  |  |  |
| FormattedMessageReceived | x |  |  |  |  |
| ImportTarget | x |  |  |  |  |
| InvalidXml |  |  |  |  |  |
| LogisticEvent | x |  |  |  |  |
| MessageReceived | x | x | x | x | x |
| OperationReceived | x |  |  |  |  |
| Position | x | x | x | x | x |
| ReceivedAreaLayout |  |  | x |  |  |
| ReceivedCheckpoint |  | x |  |  |  |
| ReceivedDetailedLayout |  |  | x |  |  |
| ReceivedDriver |  | x | x |  |  |
| ReceivedEmbarkedMacro |  | x |  |  |  |
| ReceivedGroupMacro |  | x | x |  |  |
| ReceivedItemMacro |  | x |  |  |  |
| ReceivedKeyboardLayout |  |  |  | x |  |
| ReceivedOperationalProfile |  | x |  |  |  |
| ReceivedProfileConfiguration | x |  |  |  |  |
| ReceivedRedundantVehicle |  | x |  |  |  |
| ReceivedTargetProfile | x |  |  |  |  |
| ReceivedTrailerTruck |  | x |  |  |  |
| ReceivedVehicleAccessory |  | x | x |  | x |
| ReceivedVehicle |  | x | x |  | x |
| StatusVehicle |  | x | x |  |  |
| TechnologyException | x |  |  |  |  |
| Temperature | x | x | x |  | x |
| UpdateDeviceInformation | x |  |  |  |  |
| AccountResponseSchema |  |  |  | x |  |
| AuthorizedVehicleSchema |  |  |  | x |  |
| CommandObcSchema |  |  |  | x |  |
| ExpandedAlertsSchema |  |  |  | x |  |
| MacroImportSchema |  |  |  | x |  |
| ObcProfileSchema |  |  |  | x |  |
| VehicleSchema |  | x | x | x | x |
| ReturnAuthorizationSchema |  |  |  | x |  |
| ReceivedRavexEquipmentSituationSchema |  |  |  |  | x |
| ReceivedRavexInstallationTypeSchema |  |  |  |  | x |
| ReceivedRavexMacroGroupSchema |  |  |  |  | x |
| ReceivedRavexVehicleStatusSchema |  |  |  |  | x |
| ReceivedRuleGroupSchema |  |  |  |  | x |
| ReceivedLayoutSchema |  |  | x |  |  |
| ReceiveKeyboardMacroSchema |  |  | x |  |  |
| ReceivedGroupMacroSchema (schema) |  | x | x |  |  |
| ReceivedDetailedLayoutSchema |  |  | x |  |  |
| SpySchema |  | x |  |  |  |
| ReceivedRedundantVehicleSchema |  | x |  |  |  |
| MensagemSpy |  | x |  |  |  |
| BlackBoxSchema |  | x |  |  |  |
| UpdateDeviceInformationSchema | x |  |  |  |  |
| OperationStatusSchema | x |  |  |  |  |
| BaitSchema | x |  |  |  |  |
| ApplicationEventSchema |  |  |  |  |  |
| TechnologyCommandSchema |  |  |  |  |  |
| TechnologyDataSchema |  |  |  |  |  |
| TechnologyEventSchema |  |  |  |  |  |
| TechnologyHubImportSchema |  |  |  |  |  |
| TechnologyImportSchema |  |  |  |  |  |
| TechnologyPoolingSchema |  |  |  |  |  |
| TechnologyPositionSchema |  |  |  |  |  |
| AssociateActionLayoutSchema |  |  |  |  |  |
| AssociateDriverSchema |  |  |  |  |  |
| AssociateGroupMacroSchema |  |  |  |  |  |
| AssociateLayoutSchema |  |  |  |  |  |
| AssociateOperationProfileSchema |  |  |  |  |  |
| AssociatePointGroupSchema |  |  |  |  |  |
| AssociateTargetProfileSchema |  |  |  |  |  |
| ChangePFAnalysisSchema |  |  |  |  |  |
| ChangePFSchema |  |  |  |  |  |
| ClearActionLayoutSchema |  |  |  |  |  |
| ClearDriverSchema |  |  |  |  |  |
| ClearLayoutSchema |  |  |  |  |  |
| Disable5WheelButtonSchema |  |  |  |  |  |
| Enable5WheelButtonSchema |  |  |  |  |  |
| HistoricPositionsByDateSchema |  |  |  |  |  |
| LinkTruckTrailerSchema |  |  |  |  |  |
| LockBoxSchema |  |  |  |  |  |
| MessageSchema |  |  |  |  |  |
| OperationProfileSchema |  |  |  |  |  |
| OperationSchema |  |  |  |  |  |
| PasscodeConfirmationMtc600Schema |  |  |  |  |  |
| PasscodeConfirmationSchema |  |  |  |  |  |
| RequestStatusOperationSchema |  |  |  |  |  |
| ResetSchema |  |  |  |  |  |
| SafeModeSchema |  |  |  |  |  |
| SealBoxSchema |  |  |  |  |  |
| SealEngineSchema |  |  |  |  |  |
| SealTrailerSchema |  |  |  |  |  |
| TransitionConfiguration |  |  |  |  |  |
| TransmitGroupMacroSchema |  |  |  |  |  |
| UnlinkTruckTrailerSchema |  |  |  |  |  |
| UnlockBoxSchema |  |  |  |  |  |
| UnsealBoxSchema |  |  |  |  |  |
| UnsealTrailerSchema |  |  |  |  |  |
| ImportTargetSchema |  |  |  |  |  |
| TargetGroupSchema |  |  |  |  |  |
| TargetSchema |  |  |  |  |  |
| TransmitTargetSchema |  |  |  |  |  |

> Note: Empty cells mean the schema is not listed for that technology in the original matrix.


## EventSchemaBase

Every type of packet (technology → Hub) is considered an event (positions, temperature, general events, etc.). Each event contains a standard schema and additional data depending on the event type.

Below is the standard schema:
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "json_content_type": {
      "type": "string"
    },
    "json_content": {
      "type": "string"
    },
    "account_id": {
      "type": "string",
      "format": "uuid"
    },
    "id": {
      "type": "string",
      "format": "uuid"
    },
    "technology": {
      "type": "string",
      "enum": ["Omnilink", "Sascar", "..."]
    },
    "event_time": {
      "type": "string",
      "format": "date-time"
    },
    "device_id": {
      "type": "string"
    },
    "event_type": {
      "type": "string",
      "enum": [
        "TechnologyException",
        "LogisticEvent",
        "ConfirmTransmitOperation",
        "ConfirmTransmitTargets",
        "ConfirmAssociationTargets",
        "ConfirmAssociationOperation",
        "ConfirmStatusTeleCommand",
        "UpdateDeviceInformation",
        "TargetProfile",
        "ImportationOfPredefinedMacros",
        "ConfirmReadMsgFree",
        "Temperature",
        "ArrivalTargetArea",
        "ExitTargetArea",
        "TargetArrivalCommand",
        "TargetExitCommand",
        "StartOfService",
        "EndOfService",
        "RightDoorOpening",
        "RightDoorClosing",
        "LeftDoorClosing",
        "LeftDoorOpening",
        "IgnitionOff",
        "IgnitionOn",
        "ReceivedFormattedMessage",
        "OperationReceived",
        "ReceivedTargetProfile",
        "CommandStatusReceived",
        "Authorization",
        "ReceivedPredefinedMessage",
        "InvalidXml",
        "PositionReceived",
        "ConfirmAssociationTargetsFail",
        "MessageReceive"
      ]
    }
  },
  "required": []
}
```

This is the position Schema:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "TechnologyPositionSchema",
  "description": "Schema representing a technology position.",
  "type": "object",
  "properties": {
    "id": {
      "type": "string",
      "format": "uuid",
      "description": "The position ID."
    },
    "technology_type": {
      "type": "string",
      "description": "The type of technology."
    },
    "event_time": {
      "type": "string",
      "format": "date-time",
      "description": "The event time."
    },
    "device_id": {
      "type": "string",
      "description": "The device ID."
    },
    "position_type": {
      "type": "string",
      "description": "The position type."
    },
    "latitude": {
      "type": "number",
      "description": "The latitude."
    },
    "longitude": {
      "type": "number",
      "description": "The longitude."
    },
    "ignition": {
      "type": "string",
      "description": "The ignition status."
    },
    "speed": {
      "type": "number",
      "description": "The speed."
    },
    "location": {
      "type": "string",
      "description": "The location."
    },
    "odometer": {
      "type": "integer",
      "description": "The odometer."
    }
  },
  "required": [
    "id",
    "technology_type",
    "event_time",
    "device_id",
    "position_type",
    "latitude",
    "longitude",
    "ignition",
    "speed",
    "location",
    "odometer"
  ]
}
```


### Environment viariables
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=3059"
export ELASTIC_DOC_NAME="nstech"
export NSTECH_TOKEN_URL=""
export NSTECH_URL=""
export NSTECH_USER=""
export NSTECH_USER_KEY=""
```

## Operation flow - Integra NSTech

One of Integra nstech's main characteristics is its passive interaction model — the technology initiates the process.

The first step, which happens asynchronously, is to obtain the access credentials (see HERE for the step-by-step instructions on how to obtain your credentials). This credential, which we call the Token, will be used in all operations within the platform.

The Token is reusable

It is not necessary to obtain a new token for every request; the same token used in one request can be reused while it remains valid. Request a new token only if the current one has expired.

With a valid access credential ready, the next step is to prepare the package with the information and send it through the API. We'll describe the structure of the payloads later. First, it is important to note that each type of information has a specific endpoint for submission, for example:

### API playground for CommandStatus

**Method:** POST  
**URL:** https://stg.nstech.com.br/zeus/api/integra/v1/command-status

**Request history:**  
LOG IN TO SEE FULL REQUEST HISTORY  
TIME | STATUS | USER AGENT  
Make a request to see history.  
0 Requests This Month

#### Parameters

| Name | Type | Required | Description |
|---|---:|:---:|---|
| technology_id | Guid | Yes | Identifier of the technology sending the information. The list of technologies can be obtained HERE. |
| account_id | Guid | Yes | Identifier of the client (GR, Carrier) monitoring the vehicle. The list of clients can be obtained HERE. |
| device_id | String | Yes | Tracker/device identifier. |
| date | Calendar (datetime) | Yes | Timestamp when the position was generated on the tracker. Format: yyyy-MM-ddTHH:mm:ss.fffffffZ — UTC offset: 0. |
| confirmation_id | String | Yes | Identifier sent by the GR/Carrier and made available on the technology return queue together with the command or message from a prior asynchronous process. |
| status | Integer | Yes | Execution status of the command sent by the GR to the technology. Possible values:  
- Sent — Command/message sent to the vehicle.  
- Success — Command executed successfully / Message received by the vehicle.  
- Error — Failure during execution/receipt.  
- Canceled — Sending canceled.  
- NotSupported — Command/message not supported/unknown by the technology or cannot be executed by the tracker. |

#### Query params

| Name | Type | Default |
|---|---:|---:|
| api-version | string | v1 |

#### Body params

| Name | Type |
|---|---:|
| commands | array of objects \| null |

(There is an "ADD OBJECT" UI action in the original playground to append command objects.)

#### Responses

| Status | Description |
|---:|---|
| 201 | Created |
| 400 | Bad Request |
| 422 | Client Error
#### Curl Example
```
curl --request POST \
     --url https://stg.nstech.com.br/zeus/api/integra/v1/command-status \
     --header 'accept: application/json' \
     --header 'content-type: application/*+json'
```
#### Python example:
```
import requests

url = "https://stg.nstech.com.br/zeus/api/integra/v1/command-status"

headers = {
    "accept": "application/json",
    "content-type": "application/*+json"
}

response = requests.post(url, headers=headers)

print(response.text) 
```


### API playground for Events 

**Method:** POST  
**URL:** https://stg.nstech.com.br/zeus/api/integra/v1/events

**Request history:**  
LOG IN TO SEE FULL REQUEST HISTORY  
TIME | STATUS | USER AGENT  
Make a request to see history.  
0 Requests This Month

#### Parameters

| Name | Type | Required | Description |
|---|---:|:---:|---|
| technology_id | Guid | Yes | Identifier of the technology sending the information. The list of technologies can be obtained HERE. |
| account_id | Guid | Yes | Identifier of the client (GR, Carrier) monitoring the vehicle. The list of clients can be obtained HERE. |
| device_id | Guid | Yes | Tracker/device identifier. |
| date | Calendar (datetime) | Yes | Timestamp when the position/event was generated on the tracker. Format: yyyy-MM-ddTHH:mm:ss.fffffffZ — UTC offset: 0. |
| latitude | double | Yes | Latitude in decimal format. Example: -00.000000 |
| longitude | double | Yes | Longitude in decimal format. Example: -00.000000 |
| event | string | Yes | Event code generated by the device. The list of events can be obtained HERE. |
#### Event Type table:
| Code | Name | Description |
|---|---|---|
| PanicButton | Panic Button | Event sent when the panic button is pressed. |
| BatteryConnect | Battery Connected | Event sent when the main power source is powering the vehicle. |
| BatteryDisconnect | Battery Disconnected | Event sent when the vehicle's main power source is not operating. |
| PanelViolated | Panel Violated | Event sent in case of panel tampering. |
| GPSConnected | GPS Antenna Connected | When the GPS antenna is reconnected to the tracker. |
| GPSDisconnected | GPS Antenna Disconnected | When the GPS antenna is no longer communicating due to device disconnection (not due to other signal loss reasons). |
| PotentialJammer | Potential Jammer | Communication anomalies indicating possible use of a jammer (signal blocker device). |
| BoxRearOpen | Rear Box Door Open | Sensor on the rear box door signals it has been opened. |
| BoxRearClosed | Rear Box Door Closed | Sensor on the rear box door signals it has been closed. |
| BoxRearViolated | Rear Box Door Violated | Sensor on the rear box door signals it was opened improperly. |
| BoxSideOpen | Side Box Door Open | Sensor on the side box door signals it has been opened. |
| BoxSideClosed | Side Box Door Closed | Sensor on the side box door signals it has been closed. |
| BoxSideViolated | Side Box Door Violated | Sensor on the side box door signals it was opened improperly. |
| TerminalConnected | Terminal Keypad Connected | Occurs when the tracker terminal keypad is reconnected. |
| TerminalDisconnected | Terminal Keypad Disconnected | Occurs when the tracker terminal keypad is disconnected. |
| DriveDoorOpen | Driver Door Open | Sensor on the driver door signals it has been opened. |
| DriveDoorClosed | Driver Door Closed | Sensor on the driver door signals it has been closed. |
| DriveDoorViolated | Driver Door Violated | Sensor on the driver door signals it was opened improperly. |
| PassengerDoorOpen | Passenger Door Open | Sensor on the passenger door signals it has been opened. |
| PassengerDoorClosed | Passenger Door Closed | Sensor on the passenger door signals it has been closed. |
| PassengerDoorViolated | Passenger Door Violated | Sensor on the passenger door signals it was opened improperly. |
| TrailerEngage | Trailer Engaged | Event sent when the tractor and trailer are coupled. |
| TrailerDisengage | Trailer Disengaged | Event sent when the tractor and trailer are uncoupled. |
| FiveWhellLock | Fifth Wheel Lock Activated | Sensor on the fifth wheel lock signals it has been locked. |
| FiveWhellUnlock | Fifth Wheel Lock Deactivated | Sensor on the fifth wheel lock signals it has been unlocked. |
| FiveWhellViolated | Fifth Wheel Lock Violated | Sensor on the fifth wheel lock signals it

#### Curl example:
```
curl --location 'https://stg.nstech.com.br/zeus/api/integra/v1/events?api-version=v1' \
--header 'Content-Type: application/json' \
--header 'Authorization: *******' \
--data '{
  "events": [
    {
      "technology_id": "[id da tecnologia-MUDAR]",
      "account_id": "[id da conta destino - MUDAR]",
      "date": "2024-06-18T14:34:32.882Z",
      "device_id": "868125",
      "event_type": "PanicButton",
      "latitude": -23.584492,
      "longitude": -46.828401
    }
  ]
}'
```

### API playground for Messages

**Method:** POST  
**URL:** https://stg.nstech.com.br/zeus/api/integra/v1/messages

**Request history:**  
LOG IN TO SEE FULL REQUEST HISTORY  
TIME | STATUS | USER AGENT  
Make a request to see history.  
0 Requests This Month

#### Description
Message sending process to the manager. The source of the information is the vehicle or the technology.

#### Parameters

| Name | Type | Required | Description |
|---|---:|:---:|---|
| technology_id | Guid | Yes | Identifier of the technology sending the information. The list of technologies can be obtained HERE. |
| account_id | Guid | Yes | Identifier of the client (GR, Carrier) monitoring the vehicle. The list of clients can be obtained HERE. |
| device_id | String | Yes | Tracker/device identifier. |
| date | Calendar (datetime) | Yes | Timestamp generated by the tracker. Format: yyyy-MM-ddTHH:mm:ss.fffffffZ — UTC offset: 0. |
| latitude | double | Yes | Latitude in decimal format. Example: -00.000000 |
| longitude | double | Yes | Longitude in decimal format. Example: -00.000000 |
| content | String | Yes | Message content. |

#### Query params

| Name | Type | Default |
|---|---:|---:|
| api-version | string | v1 |

#### Body params

| Name | Type |
|---|---:|
| messages | array of objects \| null |

Note: Up to 1000 objects per request.

#### Responses

| Status | Description |
|---:|---|
| 201 | Created |
| 400 | Bad Request |
| 422 | Client Error |
| 500 | Server Error |

#### Curl example:
```
curl --request POST \
     --url https://stg.nstech.com.br/zeus/api/integra/v1/messages \
     --header 'accept: application/json' \
     --header 'content-type: application/*+json'
```

### API playground for Positions

**Method:** POST  
**URL:** https://stg.nstech.com.br/zeus/api/integra/v1/positions

**Request history:**  
LOG IN TO SEE FULL REQUEST HISTORY  
TIME | STATUS | USER AGENT  
Make a request to see history.  
0 Requests This Month

## Parameters

| Name | Type | Required | Description |
|---|---:|:---:|---|
| technology_id | Guid | Yes | Identifier of the technology sending the information. The list of technologies can be obtained HERE. |
| account_id | Guid | Yes | Identifier of the client (GR, Carrier) monitoring the vehicle. The list of clients can be obtained HERE. |
| device_id | Guid | Yes | Tracker/device identifier. |
| date | Calendar (datetime) | Yes | Timestamp when the position was generated on the tracker. Format: yyyy-MM-ddTHH:mm:ss.fffffffZ — UTC offset: 0. |
| latitude | double | Yes | Latitude in decimal format. Example: -00.000000 |
| longitude | double | Yes | Longitude in decimal format. Example: -00.000000 |
| speed | decimal | No | Speed obtained from telemetry or GPS; prefer telemetry when available. |
| ignition | string | No | Ignition status. Allowed values:  
- Off  
- On  
- Unknown (default if not provided) |
| odometer | decimal | No | Distance measurement (odometer). |
| position_type | string | No | Origin of the packet:  
- GPRS — Obtained via cellular/telephony  
- Satellite

#### Curl example:
```
curl --request POST \
     --url https://stg.nstech.com.br/zeus/api/integra/v1/positions \
     --header 'accept: application/json' \
     --header 'content-type: application/*+json'
```
json content
```
{
  "positions": [
    {
      "technology_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "account_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "date": "2024-04-09T14:14:19.688Z",
      "device_id": "868125",
      "position_type": "GPRS",
      "latitude": -23.584492,
      "longitude": -46.828401,
      "ignition": "Off",
      "speed": 0,
      "odometer": 14587
    }
```

#### API playground for Temperatures

**Method:** POST  
**URL:** https://stg.nstech.com.br/zeus/api/integra/v1/temperatures

**Request history:**  
LOG IN TO SEE FULL REQUEST HISTORY  
TIME | STATUS | USER AGENT  
Make a request to see history.  
0 Requests This Month

#### Parameters

| Name | Type | Required | Description |
|---|---:|:---:|---|
| technology_id | Guid | Yes | Identifier of the technology sending the information. The list of technologies can be obtained HERE. |
| account_id | Guid | Yes | Identifier of the client (GR, Carrier) monitoring the vehicle. The list of clients can be obtained HERE. |
| device_id | String | Yes | Tracker/device identifier. |
| date | Calendar (datetime) | Yes | Timestamp when the reading was generated on the tracker. Format: yyyy-MM-ddTHH:mm:ss.fffffffZ — UTC offset: 0. |
| latitude | double | Yes | Latitude in decimal format. Example: -00.000000 |
| longitude | double | Yes | Longitude in decimal format. Example: -00.000000 |
| temperature | integer | Yes | Temperature value in degrees Celsius at the moment; can be negative. |
| sensor | integer | No | Sensor index that measured the temperature. If omitted, defaults to 1. Valid values: 1–5. |

#### Query params

| Name | Type | Default |
|---|---:|---|
| api-version | string | v1 |

#### Body params

| Name | Type |
|---|---:|
| temperatures | array of objects \| null |

Note: Up to 1000 objects per request.

#### Responses

| Status | Description |
|---:|---|
| 201 | Created |
| 400 | Bad Request |
| 422 | Client Error |
| 500 | Server Error |

After submission, the API performs an initial validation of required fields and content and quickly returns confirmation indicating whether the information was processed or not.

APIs are type-specific

Each API is prepared to receive a specific type of information. Sending data to the wrong endpoint or queue will generate an error that is reported in the request response.

Validation error responses from the requests above will be made available in an exclusive queue; those responses follow the schema shown earlier in this document.
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "type": {
      "type": "string",
      "nullable": true,
      "description": "A URI reference [RFC3986] that identifies the problem type. This specification encourages that, when dereferenced, it provide human-readable documentation for the problem type (e.g., using HTML [W3C.REC-html5-20141028]). When this member is not present, its value is assumed to be 'about:blank'."
    },
    "title": {
      "type": "string",
      "nullable": true,
      "description": "A short, human-readable summary of the problem type. It SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization(e.g., using proactive content negotiation; see[RFC7231], Section 3.4)."
    },
    "status": {
      "type": "integer",
      "nullable": true,
      "description": "The HTTP status code([RFC7231], Section 6) generated by the origin server for this occurrence of the problem."
    },
    "detail": {
      "type": "string",
      "nullable": true,
      "description": "A human-readable explanation specific to this occurrence of the problem."
    },
    "instance": {
      "type": "string",
      "nullable": true,
      "description": "A URI reference that identifies the specific occurrence of the problem. It may or may not yield further information if dereferenced."
    },
    "extensions": {
      "type": "object",
      "description": "Problem type definitions MAY extend the problem details object with additional members. Extension members appear in the same namespace as other members of a problem type.",
      "additionalProperties": {
        "type": "string"    
      }
    }
  }
}
```