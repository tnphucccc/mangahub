# MangaHub gRPC Service Documentation

Complete gRPC service documentation for the MangaHub internal service.

## Table of Contents

- [Overview](#overview)
- [Protocol Buffer Definition](#protocol-buffer-definition)
- [Service Methods](#service-methods)
- [Message Types](#message-types)
- [Client Examples](#client-examples)
- [Error Handling](#error-handling)

---

## Overview

### Service URL

```
localhost:9092
```

### Protocol

gRPC using Protocol Buffers (proto3)

### Purpose

The gRPC service provides internal service-to-service communication for:

- Manga information retrieval
- Advanced search capabilities
- Reading progress management

This service is designed for internal use by other MangaHub services and is not exposed to end users directly.

---

## Protocol Buffer Definition

The service is defined in `proto/manga.proto`:

```protobuf
syntax = "proto3";

package manga;

option go_package = "mangahub/internal/grpc/pb";

service MangaService {
    rpc GetManga(GetMangaRequest) returns (MangaResponse);
    rpc SearchManga(SearchRequest) returns (SearchResponse);
    rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
}
```

---

## Service Methods

### 1. GetManga

Retrieve detailed information about a specific manga by its ID.

**Method Signature:**

```protobuf
rpc GetManga(GetMangaRequest) returns (MangaResponse);
```

**Request Message:**

```protobuf
message GetMangaRequest {
    string manga_id = 1;
}
```

**Response Message:**

```protobuf
message MangaResponse {
    string        id             = 1;
    string        title          = 2;
    string        author         = 3;
    repeated string genres       = 4;
    string        status         = 5;
    int32         total_chapters = 6;
    string        description    = 7;
    string        cover_url      = 8;
}
```

**Status Codes:**

- `OK (0)`: Manga found and returned successfully
- `NOT_FOUND (5)`: Manga with specified ID does not exist
- `INTERNAL (13)`: Internal server error

**Example Request:**

```protobuf
{
  manga_id: "manga-001"
}
```

**Example Response:**

```protobuf
{
  id: "manga-001",
  title: "One Piece",
  author: "Eiichiro Oda",
  genres: ["Action", "Adventure", "Fantasy"],
  status: "ongoing",
  total_chapters: 1100,
  description: "The story follows Monkey D. Luffy...",
  cover_url: "https://example.com/onepiece.jpg"
}
```

---

### 2. SearchManga

Search for manga with multiple filter criteria.

**Method Signature:**

```protobuf
rpc SearchManga(SearchRequest) returns (SearchResponse);
```

**Request Message:**

```protobuf
message SearchRequest {
    string title    = 1;  // Filter by title (partial match, case-insensitive)
    string author   = 2;  // Filter by author (partial match, case-insensitive)
    string genre    = 3;  // Filter by genre
    string status   = 4;  // Filter by status (ongoing, completed, hiatus, cancelled)
    string order_by = 5;  // Order results (currently not implemented)
    int32  limit    = 6;  // Maximum results (default: 20, max: 100)
    int32  offset   = 7;  // Pagination offset (default: 0)
}
```

**Request Field Details:**

| Field      | Type   | Required | Description                                                      | Default |
| ---------- | ------ | -------- | ---------------------------------------------------------------- | ------- |
| `title`    | string | No       | Filter by manga title (case-insensitive, partial match)          | -       |
| `author`   | string | No       | Filter by author name (case-insensitive, partial match)          | -       |
| `genre`    | string | No       | Filter by genre                                                  | -       |
| `status`   | string | No       | Filter by status (`ongoing`, `completed`, `hiatus`, `cancelled`) | -       |
| `order_by` | string | No       | Order results (reserved for future use)                          | -       |
| `limit`    | int32  | No       | Number of results (max: 100)                                     | 20      |
| `offset`   | int32  | No       | Pagination offset                                                | 0       |

**Response Message:**

```protobuf
message SearchResponse {
    repeated MangaResponse manga = 1;
}
```

**Status Codes:**

- `OK (0)`: Search completed successfully (may return empty array)
- `INTERNAL (13)`: Internal server error

**Example Request (Search by Title):**

```protobuf
{
  title: "naruto",
  limit: 10,
  offset: 0
}
```

**Example Request (Search by Author):**

```protobuf
{
  author: "kishimoto",
  limit: 10,
  offset: 0
}
```

**Example Request (Combined Filters):**

```protobuf
{
  title: "one",
  genre: "Action",
  status: "ongoing",
  limit: 20,
  offset: 0
}
```

**Example Response:**

```protobuf
{
  manga: [
    {
      id: "manga-002",
      title: "Naruto",
      author: "Masashi Kishimoto",
      genres: ["Action", "Adventure", "Martial Arts"],
      status: "completed",
      total_chapters: 700,
      description: "The story follows Naruto Uzumaki...",
      cover_url: "https://example.com/naruto.jpg"
    },
    {
      id: "manga-003",
      title: "Naruto Shippuden",
      author: "Masashi Kishimoto",
      genres: ["Action", "Adventure"],
      status: "completed",
      total_chapters: 500,
      description: "Continuation of Naruto's story...",
      cover_url: "https://example.com/shippuden.jpg"
    }
  ]
}
```

---

### 3. UpdateProgress

Update a user's reading progress for a specific manga.

**Method Signature:**

```protobuf
rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
```

**Request Message:**

```protobuf
message UpdateProgressRequest {
    string user_id  = 1;  // User ID
    string manga_id = 2;  // Manga ID
    string status   = 3;  // Reading status
    int32  chapter  = 4;  // Current chapter number
    int32  rating   = 5;  // User rating (1-10)
}
```

**Request Field Details:**

| Field      | Type   | Required | Description            | Constraints                                                  |
| ---------- | ------ | -------- | ---------------------- | ------------------------------------------------------------ |
| `user_id`  | string | Yes      | User ID                | Must exist in database                                       |
| `manga_id` | string | Yes      | Manga ID               | Must exist in database                                       |
| `status`   | string | No       | Reading status         | `reading`, `completed`, `plan_to_read`, `on_hold`, `dropped` |
| `chapter`  | int32  | Yes      | Current chapter number | 0 to total_chapters                                          |
| `rating`   | int32  | No       | User rating            | 1-10                                                         |

**Reading Status Values:**

- `reading`: Currently reading
- `completed`: Finished reading
- `plan_to_read`: Planning to read
- `on_hold`: Temporarily stopped
- `dropped`: Abandoned

**Response Message:**

```protobuf
message UpdateProgressResponse {
    UserProgress progress = 1;
}

message UserProgress {
    string user_id         = 1;
    string manga_id        = 2;
    int32  current_chapter = 3;
    string status          = 4;
    int32  rating          = 5;
    string started_at      = 6;  // Unix timestamp as string
    string completed_at    = 7;  // Unix timestamp as string
    string updated_at      = 8;  // Unix timestamp as string
}
```

**Status Codes:**

- `OK (0)`: Progress updated successfully
- `INTERNAL (13)`: Internal server error (includes validation failures)

**Example Request:**

```protobuf
{
  user_id: "user-001",
  manga_id: "manga-001",
  status: "reading",
  chapter: 150,
  rating: 9
}
```

**Example Response:**

```protobuf
{
  progress: {
    user_id: "user-001",
    manga_id: "manga-001",
    current_chapter: 150,
    status: "reading",
    rating: 9,
    started_at: "1732704000",
    completed_at: "0",
    updated_at: "1732790400"
  }
}
```

---

## Message Types

### MangaResponse

Complete manga information.

```protobuf
message MangaResponse {
    string        id             = 1;  // Unique manga identifier
    string        title          = 2;  // Manga title
    string        author         = 3;  // Author name
    repeated string genres       = 4;  // Array of genre tags
    string        status         = 5;  // Publication status
    int32         total_chapters = 6;  // Total number of chapters
    string        description    = 7;  // Synopsis/description
    string        cover_url      = 8;  // Cover image URL
}
```

**Field Descriptions:**

- `id`: Unique identifier (e.g., "manga-001")
- `title`: Full manga title
- `author`: Author's name
- `genres`: Array of genre strings (e.g., ["Action", "Adventure"])
- `status`: One of: `ongoing`, `completed`, `hiatus`, `cancelled`
- `total_chapters`: Total number of published chapters
- `description`: Detailed synopsis
- `cover_url`: URL to cover image

---

### UserProgress

User's reading progress for a manga.

```protobuf
message UserProgress {
    string user_id         = 1;  // User identifier
    string manga_id        = 2;  // Manga identifier
    int32  current_chapter = 3;  // Last read chapter
    string status          = 4;  // Reading status
    int32  rating          = 5;  // User rating (1-10)
    string started_at      = 6;  // When user started (Unix timestamp)
    string completed_at    = 7;  // When user completed (Unix timestamp, "0" if not completed)
    string updated_at      = 8;  // Last update time (Unix timestamp)
}
```

**Field Descriptions:**

- `user_id`: Unique user identifier
- `manga_id`: Unique manga identifier
- `current_chapter`: Chapter number user is currently on
- `status`: Current reading status
- `rating`: User's rating (1-10, 0 if not rated)
- `started_at`: Unix timestamp when user started reading (as string)
- `completed_at`: Unix timestamp when completed (as string, "0" if not completed)
- `updated_at`: Unix timestamp of last update (as string)

---

## Client Examples

### Go Client Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "github.com/tnphucccc/mangahub/internal/grpc/pb"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewMangaServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 1. Get manga by ID
    fmt.Println("=== Get Manga by ID ===")
    getMangaResp, err := client.GetManga(ctx, &pb.GetMangaRequest{
        MangaId: "manga-001",
    })
    if err != nil {
        log.Fatalf("GetManga failed: %v", err)
    }
    fmt.Printf("Manga: %s by %s\n", getMangaResp.Title, getMangaResp.Author)
    fmt.Printf("Status: %s, Chapters: %d\n\n", getMangaResp.Status, getMangaResp.TotalChapters)

    // 2. Search manga by title
    fmt.Println("=== Search Manga by Title ===")
    searchResp, err := client.SearchManga(ctx, &pb.SearchRequest{
        Title:  "naruto",
        Limit:  10,
        Offset: 0,
    })
    if err != nil {
        log.Fatalf("SearchManga failed: %v", err)
    }
    fmt.Printf("Found %d manga:\n", len(searchResp.Manga))
    for _, m := range searchResp.Manga {
        fmt.Printf("  - %s (%s)\n", m.Title, m.Author)
    }
    fmt.Println()

    // 3. Search manga by author
    fmt.Println("=== Search Manga by Author ===")
    authorSearchResp, err := client.SearchManga(ctx, &pb.SearchRequest{
        Author: "kishimoto",
        Limit:  10,
        Offset: 0,
    })
    if err != nil {
        log.Fatalf("SearchManga failed: %v", err)
    }
    fmt.Printf("Found %d manga by this author:\n", len(authorSearchResp.Manga))
    for _, m := range authorSearchResp.Manga {
        fmt.Printf("  - %s\n", m.Title)
    }
    fmt.Println()

    // 4. Update progress
    fmt.Println("=== Update Progress ===")
    updateResp, err := client.UpdateProgress(ctx, &pb.UpdateProgressRequest{
        UserId:  "user-001",
        MangaId: "manga-001",
        Status:  "reading",
        Chapter: 150,
        Rating:  9,
    })
    if err != nil {
        log.Fatalf("UpdateProgress failed: %v", err)
    }
    fmt.Printf("Updated progress: Chapter %d, Rating: %d\n",
        updateResp.Progress.CurrentChapter,
        updateResp.Progress.Rating)
}
```

### Python Client Example

```python
import grpc
import manga_pb2
import manga_pb2_grpc

def main():
    # Connect to gRPC server
    channel = grpc.insecure_channel('localhost:9092')
    client = manga_pb2_grpc.MangaServiceStub(channel)

    # 1. Get manga by ID
    print("=== Get Manga by ID ===")
    response = client.GetManga(manga_pb2.GetMangaRequest(manga_id="manga-001"))
    print(f"Manga: {response.title} by {response.author}")
    print(f"Status: {response.status}, Chapters: {response.total_chapters}\n")

    # 2. Search manga by title
    print("=== Search Manga by Title ===")
    response = client.SearchManga(manga_pb2.SearchRequest(
        title="naruto",
        limit=10,
        offset=0
    ))
    print(f"Found {len(response.manga)} manga:")
    for m in response.manga:
        print(f"  - {m.title} ({m.author})")
    print()

    # 3. Search manga by author
    print("=== Search Manga by Author ===")
    response = client.SearchManga(manga_pb2.SearchRequest(
        author="kishimoto",
        limit=10,
        offset=0
    ))
    print(f"Found {len(response.manga)} manga by this author:")
    for m in response.manga:
        print(f"  - {m.title}")
    print()

    # 4. Update progress
    print("=== Update Progress ===")
    response = client.UpdateProgress(manga_pb2.UpdateProgressRequest(
        user_id="user-001",
        manga_id="manga-001",
        status="reading",
        chapter=150,
        rating=9
    ))
    print(f"Updated progress: Chapter {response.progress.current_chapter}, "
          f"Rating: {response.progress.rating}")

if __name__ == '__main__':
    main()
```

### Using grpcurl (Command Line)

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List available services
grpcurl -plaintext localhost:9092 list

# List methods in MangaService
grpcurl -plaintext localhost:9092 list manga.MangaService

# Get manga by ID
grpcurl -plaintext -d '{"manga_id": "manga-001"}' \
  localhost:9092 manga.MangaService/GetManga

# Search manga by title
grpcurl -plaintext -d '{"title": "naruto", "limit": 10}' \
  localhost:9092 manga.MangaService/SearchManga

# Search manga by author
grpcurl -plaintext -d '{"author": "kishimoto", "limit": 10}' \
  localhost:9092 manga.MangaService/SearchManga

# Update progress
grpcurl -plaintext -d '{
  "user_id": "user-001",
  "manga_id": "manga-001",
  "status": "reading",
  "chapter": 150,
  "rating": 9
}' localhost:9092 manga.MangaService/UpdateProgress
```

---

## Error Handling

### gRPC Status Codes

The service uses standard gRPC status codes:

| Code | Name      | Description    | When Used                              |
| ---- | --------- | -------------- | -------------------------------------- |
| 0    | OK        | Success        | Request completed successfully         |
| 5    | NOT_FOUND | Not found      | Manga or user not found                |
| 13   | INTERNAL  | Internal error | Server-side error, validation failures |

### Error Response Format

Errors are returned as gRPC status with error messages:

```go
// Example error in Go client
resp, err := client.GetManga(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        fmt.Printf("Error code: %d\n", st.Code())
        fmt.Printf("Error message: %s\n", st.Message())
    }
}
```

### Common Error Scenarios

**Manga Not Found:**

```
Code: NOT_FOUND (5)
Message: "Manga with ID manga-999 not found"
```

**Search Failed:**

```
Code: INTERNAL (13)
Message: "Failed to search manga: <details>"
```

**Update Progress Failed:**

```
Code: INTERNAL (13)
Message: "Failed to update progress: <details>"
```

---

## Testing the gRPC Service

### Starting the Server

```bash
# Start the gRPC server
go run cmd/grpc-server/main.go
```

The server will start on port 9092.

### Running Tests

```bash
# Run unit tests
go test ./internal/grpc/...

# Run with verbose output
go test -v ./internal/grpc/...

# Run integration tests
go test ./test/integration/grpc_test.go
```

### Interactive Testing

Use the included gRPC CLI client for interactive testing:

```bash
# Build the CLI
go build -o grpc-client test/grpc-client/main.go

# Run interactive commands
./grpc-client get --manga-id manga-001
./grpc-client search --title naruto
./grpc-client search --author kishimoto
./grpc-client update --user-id user-001 --manga-id manga-001 --chapter 150
```

---

## Performance Considerations

### Connection Pooling

- gRPC automatically handles connection pooling
- Reuse client connections when possible
- Use context timeouts to prevent hanging requests

### Best Practices

- Set appropriate context timeouts (5-10 seconds recommended)
- Implement retry logic with exponential backoff
- Use connection keepalive for long-lived connections
- Monitor response times and adjust limits accordingly

---

## Security Considerations

**Current Implementation:**

- No authentication required (internal service only)
- Runs on localhost only
- Not exposed to external network

**For Production:**

- Implement TLS/SSL encryption
- Add authentication (mutual TLS, API keys)
- Use service mesh for service-to-service auth
- Add rate limiting
- Implement request validation

---

## Generating Client Code

### For Go

```bash
# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/manga.proto
```

### For Python

```bash
# Install required packages
pip install grpcio grpcio-tools

# Generate Python code
python -m grpc_tools.protoc -I./proto \
  --python_out=. --grpc_python_out=. \
  proto/manga.proto
```

### For Other Languages

See [gRPC documentation](https://grpc.io/docs/languages/) for language-specific code generation instructions.

---

**Last Updated**: 2025-12-20
**API Version**: 1.0.0
**Protocol**: gRPC (10 points)
