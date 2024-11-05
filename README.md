# Take-home challenge

Goal is to create a Go-driven REST service sending messages to a chat service.

## Challenge Summary

- **Objective**: 
  - Develop a RESTful web service in Go that processes and forwards specific notifications to a Discord channel

- **Functional Requirements**:
  - The service should expose a POST endpoint to receive notifications.
  - Notifications are JSON payloads that include at least a `Type`, `Name`, and `Description`.
  - Only notifications with the `Type` set to "Warning" should be forwarded to the Discord channel.
  - Notifications with other types, such as "Info", should not be forwarded and should be ignored.

- **Technical Requirements**:
  - The service is implemented in Go.
  - The service should be designed to handle multiple incoming requests from various services.

- **Non-Functional Requirements**:
  - Keep the implementation as simple as possible; authentication is not required.
  - Use in-memory storage; no external databases should be used.
  - The solution should include unit tests to verify that the API functions as expected.
  - Tests should cover different scenarios, including valid "Warning" notifications, non-forwarded "Info" notifications, and handling of invalid payloads.


# Run the server

* Tested using Go 1.23.0
* The server forwards warnings to a discord channel using a webhook url provided via the environment variable `DISCORD_WEBHOOK_URL`. Setting up the webhook requires the following steps:
  * Install Discord
  * Create a new server
  * In the default channel click the cog
  * Navigate to integrations
  * Click "create webhook"
  * Copy the webhook URL
  * Use this URL via `DISCORD_WEBHOOK_URL=...`


```bash
DISCORD_WEBHOOK_URL=WEBHOOKURL go run cmd/server/main.go
```

## Manual testing

Using a http client like [httpie](https://github.com/httpie/cli).

Valid payloads
```bash
http POST localhost:8080/notification @payloads/backup_failure_warning.json
http POST localhost:8080/notification @payloads/quota_exceeded_info.json
```

Invalid payloads
```bash
http POST localhost:8080/notification @payloads/unknown_type.json
```

## Integration testing

* The integration test requires a webhook url provided via the environment variable `DISCORD_WEBHOOK_URL`

```bash
DISCORD_WEBHOOK_URL=WEBHOOKURL go test test/server_test.go 
```


# Implementation details

* OpenAPI Specification (OAS) is used to define and document the endpoints as well as to generate server and client related source code for the go programming language. Usually the generated code would not be commited, but to provide a fully functional example right away it is part of the tracked source code. To regenerate the OAS based code, to the following.
  ```bash
  # pull the docker iamage
  podman pull docker.io/openapitools/openapi-generator-cli

  # generate the code
  podman run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i /local/openapi_spec.yaml \
    -g go-server \
    -o /local/internal/generated/openapi \
    -c /local/generate_go_server_config.yaml
  ```
* The project follows the layout from https://github.com/golang-standards/project-layout
  ```
  ├── README.md
  ├── cmd
  │   └── server
  │       └── main.go                               # Server entry point
  ├── generate_go_server_config.yaml                # OpenAPI specific generation spec
  ├── internal
  │   ├── api                                       # api service implementation
  │   │   ├── api_notification_service.go
  │   │   └── notification_type.go
  │   ├── discord                                   # Discord client implementation to send messages via webhook
  │   │   ├── client.go
  │   │   └── config.go
  │   ├── generated                                 # OpenAPI-based generated code
  │   └── server
  │       └── server.go
  ├── openapi_spec.yaml                             # Specification for the endpoint definitions
  ├── payloads                                      # Example payloads for manual testing
  └── test                                          # Contains an integration test
  ```

# Extension points

* To add more or adjust the existing endpoints start by editing `openapi_spec.yaml` according to your needs and re-generate the API specific code (see instructions above). For editing the OAS yaml you can use the [swagger editor](https://editor.swagger.io/) for convenience.
* To support more notification types, add the types in `internal/api/notification_type.go` and complete the switch statement in `internal/api/api_notification_service.go` accordingly.
* To replace the discord specific client, implement the interface `NotificationReceiver` in a custom client and replace it in `internal/server/server.go`
* To change the message template, adjust it in `internal/api/api_notification_service.go` directly or extract a template renderer to provide more flexibility.
* The port (:8080) can be changed in `internal/server/server.go`