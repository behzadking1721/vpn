# VPN Client API Documentation

This document provides detailed information about the REST API endpoints available in the VPN Client application. The API allows you to manage servers, subscriptions, and connections programmatically.

## Base URL

All API endpoints are prefixed with `/api`. For example, to list servers, you would make a request to `/api/servers`.

## Authentication

The API does not currently require authentication. This may change in future versions for security reasons.

## Common Response Formats

### Success Response

```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### Error Response

```json
{
  "error": "Error message"
}
```

## HTTP Status Codes

The API uses standard HTTP status codes to indicate the success or failure of requests:

- `200 OK`: The request was successful
- `201 Created`: The resource was successfully created
- `204 No Content`: The request was successful but there is no content to return
- `400 Bad Request`: The request was invalid or cannot be served
- `404 Not Found`: The requested resource could not be found
- `500 Internal Server Error`: An error occurred on the server

## Rate Limiting

The API does not currently implement rate limiting. This may be added in future versions to prevent abuse.

## Server Management

### List All Servers

Returns a list of all servers in the system.

- **Endpoint**: `GET /api/servers`
- **Description**: Returns a list of all servers
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/servers
  ```
- **Response**:
  ```json
  [
    {
      "id": "server1",
      "name": "Server 1",
      "host": "server1.example.com",
      "port": 8080,
      "protocol": "vmess",
      "config": {
        "user_id": "abcd1234-abcd-1234-abcd-abcd1234abcd",
        "alter_id": 0,
        "security": "auto"
      },
      "enabled": true,
      "ping": 50,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": "server2",
      "name": "Server 2",
      "host": "server2.example.com",
      "port": 443,
      "protocol": "shadowsocks",
      "config": {
        "method": "aes-128-gcm",
        "password": "example_password"
      },
      "enabled": false,
      "ping": 120,
      "created_at": "2023-01-02T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    }
  ]
  ```

### List Enabled Servers

Returns a list of only enabled servers.

- **Endpoint**: `GET /api/servers/enabled`
- **Description**: Returns a list of only enabled servers
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/servers/enabled
  ```
- **Response**: Same as List All Servers, but filtered for enabled servers

### Get Server by ID

Returns details of a specific server by its ID.

- **Endpoint**: `GET /api/servers/{id}`
- **Description**: Returns details of a specific server
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/servers/server1
  ```
- **Response**:
  ```json
  {
    "id": "server1",
    "name": "Server 1",
    "host": "server1.example.com",
    "port": 8080,
    "protocol": "vmess",
    "config": {
      "user_id": "abcd1234-abcd-1234-abcd-abcd1234abcd",
      "alter_id": 0,
      "security": "auto"
    },
    "enabled": true,
    "ping": 50,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### Add Server

Adds a new server to the system.

- **Endpoint**: `POST /api/servers`
- **Description**: Adds a new server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/servers \
    -H "Content-Type: application/json" \
    -d '{
      "id": "server3",
      "name": "Server 3",
      "host": "server3.example.com",
      "port": 443,
      "protocol": "trojan",
      "config": {
        "password": "example_password"
      },
      "enabled": true
    }'
  ```
- **Request Body**:
  ```json
  {
    "id": "server3",
    "name": "Server 3",
    "host": "server3.example.com",
    "port": 443,
    "protocol": "trojan",
    "config": {
      "password": "example_password"
    },
    "enabled": true
  }
  ```
- **Response**: Returns the created server object

### Update Server

Updates an existing server.

- **Endpoint**: `PUT /api/servers/{id}`
- **Description**: Updates an existing server
- **Example Request**:
  ```bash
  curl -X PUT http://localhost:8080/api/servers/server1 \
    -H "Content-Type: application/json" \
    -d '{
      "id": "server1",
      "name": "Updated Server 1",
      "host": "server1.example.com",
      "port": 8080,
      "protocol": "vmess",
      "config": {
        "user_id": "abcd1234-abcd-1234-abcd-abcd1234abcd",
        "alter_id": 0,
        "security": "auto"
      },
      "enabled": true
    }'
  ```
- **Request Body**: Same as Add Server
- **Response**: Returns the updated server object

### Delete Server

Deletes a server from the system.

- **Endpoint**: `DELETE /api/servers/{id}`
- **Description**: Deletes a server
- **Example Request**:
  ```bash
  curl -X DELETE http://localhost:8080/api/servers/server1
  ```
- **Response**: `204 No Content`

### Enable Server

Enables a server for use.

- **Endpoint**: `POST /api/servers/{id}/enable`
- **Description**: Enables a server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/servers/server2/enable
  ```
- **Response**:
  ```json
  {
    "status": "enabled"
  }
  ```

### Disable Server

Disables a server.

- **Endpoint**: `POST /api/servers/{id}/disable`
- **Description**: Disables a server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/servers/server1/disable
  ```
- **Response**:
  ```json
  {
    "status": "disabled"
  }
  ```

### Update Server Ping

Updates the ping value for a server.

- **Endpoint**: `PUT /api/servers/{id}/ping`
- **Description**: Updates the ping value for a server
- **Example Request**:
  ```bash
  curl -X PUT http://localhost:8080/api/servers/server1/ping \
    -H "Content-Type: application/json" \
    -d '{"ping": 45}'
  ```
- **Request Body**:
  ```json
  {
    "ping": 45
  }
  ```
- **Response**:
  ```json
  {
    "ping": 45
  }
  ```

### Test Server Ping

Tests the ping for a specific server.

- **Endpoint**: `POST /api/servers/{id}/test-ping`
- **Description**: Tests the ping for a specific server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/servers/server1/test-ping
  ```
- **Response**:
  ```json
  {
    "ping": 45
  }
  ```

### Test All Servers Ping

Tests the ping for all enabled servers.

- **Endpoint**: `POST /api/servers/test-all-ping`
- **Description**: Tests the ping for all enabled servers
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/servers/test-all-ping
  ```
- **Response**:
  ```json
  [
    {
      "server_id": "server1",
      "ping": 45
    },
    {
      "server_id": "server2",
      "ping": 120
    }
  ]
  ```

### Get Best Server

Finds and returns the best server based on comprehensive testing.

- **Endpoint**: `GET /api/servers/best`
- **Description**: Finds and returns the best server based on comprehensive testing
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/servers/best
  ```
- **Response**: Returns the best server object

## Subscription Management

### List All Subscriptions

Returns a list of all subscriptions.

- **Endpoint**: `GET /api/subscriptions`
- **Description**: Returns a list of all subscriptions
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/subscriptions
  ```
- **Response**:
  ```json
  [
    {
      "id": "sub1",
      "name": "Subscription 1",
      "url": "https://example.com/sub",
      "auto_update": true,
      "server_count": 10,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "last_update": "2023-01-01T00:00:00Z"
    }
  ]
  ```

### Get Subscription by ID

Returns details of a specific subscription by its ID.

- **Endpoint**: `GET /api/subscriptions/{id}`
- **Description**: Returns details of a specific subscription
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/subscriptions/sub1
  ```
- **Response**: Same as List All Subscriptions, but for a single subscription

### Add Subscription

Adds a new subscription and imports servers from it.

- **Endpoint**: `POST /api/subscriptions`
- **Description**: Adds a new subscription and imports servers from it
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/subscriptions \
    -H "Content-Type: application/json" \
    -d '{
      "id": "sub2",
      "name": "Subscription 2",
      "url": "https://example2.com/sub",
      "auto_update": false
    }'
  ```
- **Request Body**:
  ```json
  {
    "id": "sub2",
    "name": "Subscription 2",
    "url": "https://example2.com/sub",
    "auto_update": false
  }
  ```
- **Response**: Returns the created subscription object

### Update Subscription

Updates an existing subscription.

- **Endpoint**: `PUT /api/subscriptions/{id}`
- **Description**: Updates an existing subscription
- **Example Request**:
  ```bash
  curl -X PUT http://localhost:8080/api/subscriptions/sub1 \
    -H "Content-Type: application/json" \
    -d '{
      "id": "sub1",
      "name": "Updated Subscription 1",
      "url": "https://example.com/sub",
      "auto_update": true
    }'
  ```
- **Request Body**: Same as Add Subscription
- **Response**: Returns the updated subscription object

### Delete Subscription

Deletes a subscription.

- **Endpoint**: `DELETE /api/subscriptions/{id}`
- **Description**: Deletes a subscription
- **Example Request**:
  ```bash
  curl -X DELETE http://localhost:8080/api/subscriptions/sub1
  ```
- **Response**: `204 No Content`

### Update Subscription Servers

Updates servers from a subscription.

- **Endpoint**: `POST /api/subscriptions/{id}/update`
- **Description**: Updates servers from a subscription
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/subscriptions/sub1/update
  ```
- **Response**: Returns the updated subscription object

## Connection Management

### Connect to Server

Connects to a specific server.

- **Endpoint**: `POST /api/connect`
- **Description**: Connects to a specific server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/connect \
    -H "Content-Type: application/json" \
    -d '{"server_id": "server1"}'
  ```
- **Request Body**:
  ```json
  {
    "server_id": "server1"
  }
  ```
- **Response**:
  ```json
  {
    "status": "connected",
    "server_id": "server1",
    "server": "Server 1"
  }
  ```

### Connect to Fastest Server

Connects to the server with the fastest ping.

- **Endpoint**: `POST /api/connect/fastest`
- **Description**: Connects to the server with the fastest ping
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/connect/fastest
  ```
- **Response**:
  ```json
  {
    "status": "connected",
    "server_id": "server1",
    "server": "Server 1"
  }
  ```

### Connect to Best Server

Connects to the best server based on comprehensive testing.

- **Endpoint**: `POST /api/connect/best`
- **Description**: Connects to the best server based on comprehensive testing
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/connect/best
  ```
- **Response**:
  ```json
  {
    "status": "connected",
    "server_id": "server1",
    "server": "Server 1"
  }
  ```

### Disconnect

Disconnects from the current server.

- **Endpoint**: `POST /api/disconnect`
- **Description**: Disconnects from the current server
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8080/api/disconnect
  ```
- **Response**:
  ```json
  {
    "status": "disconnected"
  }
  ```

### Get Connection Status

Returns the current connection status.

- **Endpoint**: `GET /api/status`
- **Description**: Returns the current connection status
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/status
  ```
- **Response**:
  ```json
  {
    "status": "Connected",
    "status_code": 2
  }
  ```

### Get Connection Statistics

Returns connection statistics.

- **Endpoint**: `GET /api/stats`
- **Description**: Returns connection statistics
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/api/stats
  ```
- **Response**:
  ```json
  {
    "connected": true,
    "uptime": 3600000000000,
    "data_sent": 1024,
    "data_recv": 2048,
    "server_id": "server1"
  }
  ```

## Health Check

### Check Server Health

Returns server health status.

- **Endpoint**: `GET /health`
- **Description**: Returns server health status
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8080/health
  ```
- **Response**:
  ```json
  {
    "status": "ok"
  }
  ```
  This API documentation provides a comprehensive overview of the available endpoints and their corresponding functionality. Each endpoint is described in detail, including its endpoint, description, example request, request body (if applicable), response, and response body (if applicable).
   
 