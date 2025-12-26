# WebSocket Chat Server Documentation

---

## 1. Overview

The WebSocket Chat Server is a core component of MangaHub, providing real-time chat functionality for manga discussions. It uses the WebSocket protocol for bidirectional, persistent connections enabling instant message delivery.

This server implements a Hub-based architecture pattern for efficient message broadcasting to multiple clients in different chat rooms. Users can join room-specific discussions or participate in the general chat.

**Key Features:**

- **Multi-Room Support:** Users can join different chat rooms for specific manga discussions
- **Real-Time Messaging:** Instant message delivery to all participants in a room
- **User Presence:** Join/leave notifications when users enter or exit rooms
- **Persistent Connections:** WebSocket connections remain open for low-latency communication
- **Hub-Based Architecture:** Efficient message broadcasting using Go channels

**Protocol Details:**

- **Port:** `9093`
- **Protocol:** WebSocket (HTTP Upgrade)
- **Data Format:** JSON
- **Endpoint:** `ws://localhost:9093/ws`

---

## 2. Architecture

### Hub Pattern

The WebSocket server uses a centralized **Hub** pattern:

```
┌─────────────────────────────────────────────┐
│              WebSocket Hub                  │
│                                             │
│  ┌─────────────────────────────────────┐    │
│  │  Rooms Map                          │    │
│  │  ┌──────────┐  ┌──────────┐         │    │
│  │  │ general  │  │one-piece │  ...    │    │
│  │  │ [clients]│  │ [clients]│         │    │
│  │  └──────────┘  └──────────┘         │    │
│  └─────────────────────────────────────┘    │
│                                             │
│  Channels:                                  │
│  • Broadcast  → Messages to rooms           │
│  • Register   → New client connections      │
│  • Unregister → Client disconnections       │
└─────────────────────────────────────────────┘
         ↓                    ↑
    [Clients]            [Clients]
```

### Components

1. **Hub** (`hub.go`): Central message router managing all rooms and clients
2. **Client** (`client.go`): Represents a WebSocket connection with read/write pumps
3. **Handler** (`handler.go`): HTTP handler for WebSocket upgrade
4. **Message Models** (`pkg/models/websocket_message.go`): Message type definitions

---

## 3. Client-Server Interaction Flow

A typical client session involves:

1. **HTTP Upgrade Request**
   - Client sends HTTP request to `/ws?username=johndoe&room=general`
   - Server upgrades connection from HTTP to WebSocket protocol
   - New Client object created and registered with Hub

2. **Hub Registration**
   - Client automatically joins specified room (defaults to "general")
   - Hub sends join notification to all users in that room
   - Client starts read and write goroutines (pumps)

3. **Message Exchange**
   - **Client → Server**: User sends chat messages through WebSocket
   - **Server → Hub**: Messages forwarded to Hub's broadcast channel
   - **Hub → Clients**: Hub broadcasts to all clients in target room

4. **Room Management**
   - Clients can join different rooms by sending `join_room` messages
   - Leaving a room sends `leave_room` message
   - System automatically sends join/leave notifications

5. **Disconnection**
   - Client closes connection or network fails
   - Client automatically unregistered from Hub
   - Leave notification sent to room participants
   - Connection resources cleaned up

---

## 4. Message Types

All communication uses a unified JSON message structure differentiated by the `type` field.

| Type         | Direction       | Description                              |
| ------------ | --------------- | ---------------------------------------- |
| `join_room`  | Client → Server | Request to join a specific chat room     |
| `leave_room` | Client → Server | Request to leave current room            |
| `chat`       | Bidirectional   | Regular chat message                     |
| `system`     | Server → Client | System notifications (join/leave events) |
| `error`      | Server → Client | Error messages from server               |

---

## 5. JSON Message Payloads

### Base Message Structure

```json
{
  "type": "string (WebSocketMessageType)",
  "username": "string (optional)",
  "room": "string (optional)",
  "content": "string",
  "timestamp": "RFC3339 timestamp"
}
```

### `join_room` (Client → Server)

Request to join a specific chat room:

```json
{
  "type": "join_room",
  "username": "johndoe",
  "room": "one-piece",
  "timestamp": "2025-12-26T10:00:00Z"
}
```

**Server Response**: System message broadcast to room

```json
{
  "type": "system",
  "room": "one-piece",
  "content": "johndoe joined the chat",
  "timestamp": "2025-12-26T10:00:01Z"
}
```

### `leave_room` (Client → Server)

Request to leave current room:

```json
{
  "type": "leave_room",
  "username": "johndoe",
  "room": "one-piece",
  "timestamp": "2025-12-26T10:05:00Z"
}
```

**Server Response**: System message broadcast to room

```json
{
  "type": "system",
  "room": "one-piece",
  "content": "johndoe left the chat",
  "timestamp": "2025-12-26T10:05:01Z"
}
```

### `chat` (Client → Server → Clients)

Regular chat message:

**From Client:**

```json
{
  "type": "chat",
  "username": "johndoe",
  "room": "one-piece",
  "content": "Chapter 1150 was amazing!",
  "timestamp": "2025-12-26T10:02:30Z"
}
```

**Broadcast to Room:** (Same format, delivered to all clients in room)

### `system` (Server → Client)

System notifications (auto-generated):

```json
{
  "type": "system",
  "room": "general",
  "content": "alice joined the chat",
  "timestamp": "2025-12-26T10:01:15Z"
}
```

### `error` (Server → Client)

Error message from server:

```json
{
  "type": "error",
  "content": "Room 'invalid-room' does not exist",
  "timestamp": "2025-12-26T10:03:00Z"
}
```

---

## 6. Connection Parameters

### WebSocket Endpoint

```
ws://localhost:9093/ws?username={username}&room={room}
```

**Query Parameters:**

- `username` (required): Unique username for the chat session
- `room` (optional): Initial room to join (defaults to "general")

**Example:**

```javascript
const ws = new WebSocket(
  "ws://localhost:9093/ws?username=johndoe&room=one-piece"
);
```

---

## 7. Client Implementation Example

### JavaScript/Browser Client

```javascript
// Connect to WebSocket server
const username = "johndoe";
const room = "one-piece";
const ws = new WebSocket(
  `ws://localhost:9093/ws?username=${username}&room=${room}`
);

// Connection opened
ws.onopen = (event) => {
  console.log("Connected to chat server");
};

// Receive messages
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case "chat":
      console.log(`[${message.username}]: ${message.content}`);
      break;
    case "system":
      console.log(`[SYSTEM]: ${message.content}`);
      break;
    case "error":
      console.error(`[ERROR]: ${message.content}`);
      break;
  }
};

// Send chat message
function sendMessage(content) {
  const message = {
    type: "chat",
    username: username,
    room: room,
    content: content,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(message));
}

// Join different room
function joinRoom(newRoom) {
  const message = {
    type: "join_room",
    username: username,
    room: newRoom,
    timestamp: new Date().toISOString(),
  };
  ws.send(JSON.stringify(message));
}

// Close connection
ws.onclose = (event) => {
  console.log("Disconnected from chat server");
};

// Error handling
ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
```

### Go Client Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/gorilla/websocket"
    "github.com/tnphucccc/mangahub/pkg/models"
)

func main() {
    // Connect to server
    url := "ws://localhost:9093/ws?username=goclient&room=general"
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer conn.Close()

    // Read messages
    go func() {
        for {
            var msg models.WebSocketMessage
            err := conn.ReadJSON(&msg)
            if err != nil {
                log.Println("Read error:", err)
                return
            }

            fmt.Printf("[%s] %s: %s\n", msg.Type, msg.Username, msg.Content)
        }
    }()

    // Send a message
    chatMsg := models.NewChatMessage("goclient", "general", "Hello from Go!")
    if err := conn.WriteJSON(chatMsg); err != nil {
        log.Fatal("Write error:", err)
    }

    // Keep connection alive
    select {}
}
```

---

## 8. Room Management

### Default Rooms

- **general**: Default room for general manga discussions
- **{manga-id}**: Dynamic rooms created for specific manga (e.g., "one-piece", "naruto")

### Room Lifecycle

1. **Creation**: Rooms are created automatically when the first user joins
2. **Active**: Room exists as long as at least one user is present
3. **Cleanup**: Empty rooms are automatically deleted when the last user leaves

### Switching Rooms

Clients can switch between rooms by sending `join_room` messages. The server handles:

- Removing client from previous room
- Adding client to new room
- Sending leave notification to old room
- Sending join notification to new room

---

## 9. Testing

### Manual Testing with `wscat`

Install `wscat` (WebSocket CLI client):

```bash
npm install -g wscat
```

Connect and test:

```bash
# Connect to server
wscat -c "ws://localhost:9093/ws?username=testuser&room=general"

# Send a chat message (type in terminal)
{"type":"chat","username":"testuser","room":"general","content":"Hello!","timestamp":"2025-12-26T10:00:00Z"}

# Join different room
{"type":"join_room","username":"testuser","room":"one-piece","timestamp":"2025-12-26T10:01:00Z"}
```

### Automated Testing

The repository includes WebSocket test files in `test/` directory:

```bash
# Start WebSocket server
make run-websocket

# Run tests (in separate terminal)
go test ./internal/websocket/... -v
```

---

## 10. Performance Considerations

### Concurrent Connections

- **Design Goal**: Support 50-100 concurrent users (per project spec)
- **Architecture**: Each client has dedicated read/write goroutines
- **Message Buffering**: Client send channels have buffer to prevent blocking

### Resource Management

- **Memory**: Active connections consume ~100KB per client
- **Goroutines**: 2 goroutines per connected client (read pump + write pump)
- **Channels**: Buffered channels prevent deadlocks on slow clients

### Best Practices

1. **Ping/Pong**: Clients should implement ping/pong for connection health
2. **Heartbeat**: Send periodic messages to keep connection alive
3. **Reconnection**: Implement automatic reconnection on disconnect
4. **Buffering**: Buffer messages client-side during temporary disconnections

---

## 11. Security Considerations

### Authentication

- Currently uses query parameter username (suitable for academic project)
- Production systems should use JWT tokens for authentication
- Consider implementing user verification before allowing joins

### Rate Limiting

- No rate limiting implemented (academic project scope)
- Production should limit messages per user per time window

### Input Validation

- Server validates message structure
- Content length should be limited to prevent abuse
- Consider implementing profanity filters for production

---

## 12. Troubleshooting

### Common Issues

**Connection Refused**

```
Error: dial tcp [::1]:9093: connect: connection refused
```

**Solution**: Ensure WebSocket server is running on port 9093

**Messages Not Received**

- Check that both sender and receiver are in the same room
- Verify message format matches expected JSON structure
- Check browser console for WebSocket errors

**Connection Drops**

- Implement ping/pong heartbeat mechanism
- Check network stability
- Verify firewall settings allow WebSocket connections

### Debug Logging

Server logs show:

- User join/leave events with timestamps
- Room creation/deletion
- Message broadcast operations
- Connection errors

---

## 13. Integration with MangaHub

### HTTP API Integration

WebSocket server runs alongside HTTP API server:

- HTTP API handles user authentication
- WebSocket uses authenticated sessions
- Both servers share same database

### TCP/UDP Coordination

- WebSocket handles real-time chat
- TCP handles progress synchronization
- UDP handles notifications
- All protocols can coexist and complement each other

---

## 14. Future Enhancements

Potential improvements (beyond academic scope):

1. **Private Messaging**: Direct messages between users
2. **Message History**: Persist chat history to database
3. **File Sharing**: Share images/links in chat
4. **Typing Indicators**: Show when users are typing
5. **Read Receipts**: Track message read status
6. **User Roles**: Moderators, admins with special permissions
7. **Message Reactions**: Emoji reactions to messages
8. **Search**: Search through chat history

---

## 15. References

- **Protocol Specification**: See `mangahub_project_spec.pdf` (Page 5, Section 4)
- **Use Cases**: See `mangahub_usecase (reference).pdf` (Pages 6-7, UC-011 to UC-013)
- **Implementation**: `internal/websocket/` directory
- **Models**: `pkg/models/websocket_message.go`
- **WebSocket RFC**: [RFC 6455](https://datatracker.ietf.org/doc/html/rfc6455)
- **Gorilla WebSocket**: [Documentation](https://github.com/gorilla/websocket)

---

**Last Updated**: 2025-12-26
**Version**: 1.0.0
**Status**: ✅ Fully Implemented
