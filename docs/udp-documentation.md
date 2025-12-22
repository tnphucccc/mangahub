# UDP Notification Server Documentation

---

## 1. Overview

The UDP Notification Server is a core component of MangaHub, responsible for sending real-time notifications to clients about new manga chapter releases. It uses the UDP protocol for low-latency, "fire-and-forget" message broadcasting.

This server is designed to handle a large number of concurrent clients with minimal overhead. It is a one-way notification system; clients register to listen for broadcasts but cannot send data back other than for registration and heartbeat purposes.

**Key Features:**

-   **Client Registration:** Clients can register to receive notifications.
-   **Notification Broadcasting:** Broadcasts new chapter release information to all registered clients.
-   **Heartbeat Mechanism:** A ping/pong mechanism ensures clients are still active.
-   **Automatic Cleanup:** Stale clients that have not sent a ping recently are automatically removed.

**Protocol Details:**

-   **Port:** `9091`
-   **Protocol:** UDP
-   **Data Format:** JSON

---

## 2. Client-Server Interaction Flow

A typical client session involves the following steps:

1.  **Client Starts**: The client application starts.
2.  **Send Registration**: The client sends a `register` message to the server, providing a unique `client_id`.
3.  **Receive Confirmation**:
    -   The server receives the registration and adds the client's address to its list of registered clients.
    -   It responds with a `register_success` message.
    -   If registration fails (e.g., missing `client_id`), the server sends a `register_failed` message.
4.  **Heartbeat (Ping/Pong)**:
    -   Periodically (e.g., every 30-60 seconds), the client should send a `ping` message to the server.
    -   The server immediately responds with a `pong` message.
    -   This updates the client's `LastSeen` timestamp on the server, keeping it from being marked as stale.
5.  **Listen for Notifications**:
    -   The client continuously listens for incoming UDP packets.
    -   When another part of the MangaHub system (e.g., the API server) triggers a new chapter release, the UDP server broadcasts a `notification` message to all registered clients.
6.  **Unregister on Exit**:
    -   When the client application is shutting down, it should send an `unregister` message to the server.
    -   The server removes the client from its list of registered clients to stop sending it notifications.

---

## 3. Message Types

All communication is done via a single JSON object structure, `UDPMessage`, differentiated by the `type` field.

| Type                    | Direction         | Description                                        |
| ----------------------- | ----------------- | -------------------------------------------------- |
| `register`              | Client → Server   | Register to receive notifications.                 |
| `unregister`            | Client → Server   | Unregister and stop receiving notifications.       |
| `ping`                  | Client → Server   | Heartbeat message to keep the connection alive.    |
| `register_success`      | Server → Client   | Confirms successful registration.                  |
| `register_failed`       | Server → Client   | Indicates that registration failed.                |
| `pong`                  | Server → Client   | The response to a client's `ping`.                 |
| `notification`          | Server → Client   | A broadcast message about a new chapter release.   |
| `error`                 | Server → Client   | A generic error message.                           |

---

## 4. JSON Message Payloads

### Base Message Structure

```json
{
  "type": "string (UDPMessageType)",
  "timestamp": "string (RFC3339)",
  "data": {}
}
```

### `register` (Client → Server)

-   **Type:** `register`
-   **Data:** `UDPRegisterMessage`

```json
{
  "type": "register",
  "timestamp": "2025-12-22T10:00:00Z",
  "data": {
    "client_id": "unique-client-identifier-123",
    "user_id": "user-uuid-456",
    "username": "johndoe"
  }
}
```

### `register_success` (Server → Client)

-   **Type:** `register_success`
-   **Data:** `UDPRegisterSuccessMessage`

```json
{
  "type": "register_success",
  "timestamp": "2025-12-22T10:00:01Z",
  "data": {
    "client_id": "unique-client-identifier-123",
    "message": "Registration successful"
  }
}
```

### `unregister` (Client → Server)

-   **Type:** `unregister`
-   **Data:** `UDPUnregisterMessage`

```json
{
  "type": "unregister",
  "timestamp": "2025-12-22T10:05:00Z",
  "data": {
    "client_id": "unique-client-identifier-123"
  }
}
```

### `ping` (Client → Server)

-   **Type:** `ping`
-   **Data:** `UDPPingMessage`

```json
{
  "type": "ping",
  "timestamp": "2025-12-22T10:01:00Z",
  "data": {
    "client_time": "2025-12-22T10:01:00Z"
  }
}
```

### `pong` (Server → Client)

-   **Type:** `pong`
-   **Data:** `UDPPongMessage`

```json
{
  "type": "pong",
  "timestamp": "2025-12-22T10:01:01Z",
  "data": {
    "server_time": "2025-12-22T10:01:01Z"
  }
}
```

### `notification` (Server → Client)

-   **Type:** `notification`
-   **Data:** `UDPNotification`

```json
{
  "type": "notification",
  "timestamp": "2025-12-22T11:00:00Z",
  "data": {
    "manga_id": "one-piece",
    "manga_title": "One Piece",
    "chapter_number": 1150,
    "chapter_title": "The Dawn of a New Era",
    "release_date": "2025-12-22T10:59:00Z",
    "message": "A new chapter of One Piece has been released!"
  }
}
```

---

## 5. Testing

A simple test client is available in the repository to test the UDP server's functionality.

1.  **Start the UDP Server**
    In one terminal, run the server:
    ```bash
    make run-udp
    ```
    You should see the log: `UDP Notification Server listening on :9091`

2.  **Run the Test Client**
    In a second terminal, run the simple test client:
    ```bash
    go run test/udp-simple/main.go
    ```

The client will automatically perform the following actions and print the results:
- Connect to the server.
- Register itself as `test-client-1`.
- Send a `ping` and wait for a `pong`.
- Listen for notifications for 30 seconds.
- Unregister itself before exiting.

You can observe the corresponding log messages (client registered, unregistered, etc.) in the UDP server's terminal window.
