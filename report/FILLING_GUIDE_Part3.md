# MangaHub Report - Step-by-Step Content Filling Guide (Part 3)

This is Part 3 of the content filling guide, completing Chapter 4 and covering Chapters 5-6, References, and Appendices.

---

## CHAPTER 4: IMPLEMENTATION (Continued)

### [UDP Server Introduction]
```
The UDP notification server broadcasts chapter release announcements to subscribed clients using the User Datagram Protocol. UDP's connectionless nature makes it ideal for one-to-many broadcasting where occasional message loss is acceptable. The server maintains a subscription registry mapping manga IDs to client addresses and sends notifications via UDP datagrams.
```

### [Server Setup Code]
```go
package main

import (
    "encoding/json"
    "log"
    "net"
    "sync"
    "time"
)

type Subscription struct {
    UserID    string
    MangaIDs  []string
    Address   *net.UDPAddr
    LastSeen  time.Time
}

type NotificationServer struct {
    conn          *net.UDPConn
    subscriptions map[string]*Subscription  // userID -> subscription
    mu            sync.RWMutex
}

func NewNotificationServer(port string) (*NotificationServer, error) {
    addr, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        return nil, err
    }
    
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return nil, err
    }
    
    return &NotificationServer{
        conn:          conn,
        subscriptions: make(map[string]*Subscription),
    }, nil
}

func (s *NotificationServer) Start() {
    log.Printf("UDP Notification server listening on %s", s.conn.LocalAddr())
    
    // Start cleanup routine for stale subscriptions
    go s.cleanupRoutine()
    
    // Start notification worker
    go s.notificationWorker()
    
    // Handle incoming subscription requests
    buffer := make([]byte, 4096)
    for {
        n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
        if err != nil {
            log.Printf("Read error: %v", err)
            continue
        }
        
        go s.handleMessage(buffer[:n], remoteAddr)
    }
}

func main() {
    server, err := NewNotificationServer(":9091")
    if err != nil {
        log.Fatal("Failed to start UDP server:", err)
    }
    
    server.Start()
}
```

### [Subscription Management]
```go
type SubscriptionRequest struct {
    Type     string   `json:"type"`
    UserID   string   `json:"user_id"`
    MangaIDs []string `json:"manga_ids"`
}

func (s *NotificationServer) handleMessage(data []byte, addr *net.UDPAddr) {
    var req SubscriptionRequest
    if err := json.Unmarshal(data, &req); err != nil {
        log.Printf("Invalid message format: %v", err)
        return
    }
    
    switch req.Type {
    case "subscribe":
        s.subscribe(req.UserID, req.MangaIDs, addr)
    case "unsubscribe":
        s.unsubscribe(req.UserID, req.MangaIDs)
    default:
        log.Printf("Unknown message type: %s", req.Type)
    }
}

func (s *NotificationServer) subscribe(userID string, mangaIDs []string, addr *net.UDPAddr) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Update or create subscription
    if sub, exists := s.subscriptions[userID]; exists {
        // Add new manga IDs to existing subscription
        sub.MangaIDs = uniqueMerge(sub.MangaIDs, mangaIDs)
        sub.Address = addr
        sub.LastSeen = time.Now()
    } else {
        s.subscriptions[userID] = &Subscription{
            UserID:   userID,
            MangaIDs: mangaIDs,
            Address:  addr,
            LastSeen: time.Now(),
        }
    }
    
    log.Printf("User %s subscribed to %d manga", userID, len(mangaIDs))
}

func (s *NotificationServer) unsubscribe(userID string, mangaIDs []string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    sub, exists := s.subscriptions[userID]
    if !exists {
        return
    }
    
    // Remove specific manga IDs
    sub.MangaIDs = removeMangaIDs(sub.MangaIDs, mangaIDs)
    
    // Delete subscription if no manga remaining
    if len(sub.MangaIDs) == 0 {
        delete(s.subscriptions, userID)
    }
    
    log.Printf("User %s unsubscribed from %d manga", userID, len(mangaIDs))
}

func (s *NotificationServer) cleanupRoutine() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        s.mu.Lock()
        cutoff := time.Now().Add(-15 * time.Minute)
        
        for userID, sub := range s.subscriptions {
            if sub.LastSeen.Before(cutoff) {
                delete(s.subscriptions, userID)
                log.Printf("Cleaned up stale subscription: %s", userID)
            }
        }
        s.mu.Unlock()
    }
}
```

### [Broadcasting Logic]
```go
type ChapterNotification struct {
    Type        string    `json:"type"`
    MangaID     string    `json:"manga_id"`
    MangaTitle  string    `json:"manga_title"`
    Chapter     int       `json:"chapter"`
    ReleaseDate time.Time `json:"release_date"`
}

func (s *NotificationServer) notificationWorker() {
    // Poll database for new chapters every minute
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    lastCheck := time.Now()
    
    for range ticker.C {
        // Query database for chapters released since last check
        newChapters := getNewChapters(lastCheck)
        lastCheck = time.Now()
        
        for _, chapter := range newChapters {
            s.broadcastChapterRelease(chapter)
        }
    }
}

func (s *NotificationServer) broadcastChapterRelease(chapter ChapterNotification) {
    chapter.Type = "chapter_release"
    
    data, err := json.Marshal(chapter)
    if err != nil {
        log.Printf("Failed to marshal notification: %v", err)
        return
    }
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    sent := 0
    for _, sub := range s.subscriptions {
        // Check if user subscribed to this manga
        if contains(sub.MangaIDs, chapter.MangaID) {
            _, err := s.conn.WriteToUDP(data, sub.Address)
            if err != nil {
                log.Printf("Failed to send notification to %s: %v", sub.UserID, err)
            } else {
                sent++
            }
        }
    }
    
    log.Printf("Broadcast chapter %d of %s to %d subscribers", 
        chapter.Chapter, chapter.MangaTitle, sent)
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### [WebSocket Introduction]
```
The WebSocket chat system provides real-time bidirectional communication for manga discussions. Built using the Gorilla WebSocket library, the implementation handles HTTP-to-WebSocket upgrades, manages multiple chat rooms, and broadcasts messages efficiently. The hub pattern coordinates message distribution across connected clients within each room.
```

### [Connection Upgrade Code]
```go
package chat

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // In production, validate origin properly
        return true
    },
}

type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    username string
    room     string
}

type Message struct {
    Type      string    `json:"type"`
    Room      string    `json:"room"`
    Username  string    `json:"username"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
    UserCount int       `json:"user_count,omitempty"`
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    
    // Extract username from query parameter or JWT
    username := r.URL.Query().Get("username")
    if username == "" {
        username = "Anonymous"
    }
    
    // Create client
    client := &Client{
        hub:      hub,
        conn:     conn,
        send:     make(chan []byte, 256),
        username: username,
        room:     "",  // Will be set when joining room
    }
    
    // Start read and write pumps
    go client.writePump()
    go client.readPump()
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        _, data, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        var msg Message
        if err := json.Unmarshal(data, &msg); err != nil {
            log.Printf("Invalid message format: %v", err)
            continue
        }
        
        msg.Username = c.username  // Ensure username matches client
        msg.Timestamp = time.Now()
        
        switch msg.Type {
        case "join":
            c.joinRoom(msg.Room)
        case "message":
            if c.room != "" {
                c.hub.broadcast <- &msg
            }
        case "typing":
            if c.room != "" {
                c.hub.broadcast <- &msg
            }
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)
            
            // Add queued messages to current websocket message
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }
            
            if err := w.Close(); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) joinRoom(room string) {
    // Leave current room if in one
    if c.room != "" {
        c.hub.leaveRoom <- c
    }
    
    // Join new room
    c.room = room
    c.hub.joinRoom <- c
    
    // Notify room of new user
    msg := &Message{
        Type:     "user_joined",
        Room:     room,
        Username: c.username,
    }
    c.hub.broadcast <- msg
}
```

### [Hub Pattern]
```go
type Hub struct {
    rooms      map[string]map[*Client]bool  // room -> clients
    broadcast  chan *Message
    joinRoom   chan *Client
    leaveRoom  chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        rooms:      make(map[string]map[*Client]bool),
        broadcast:  make(chan *Message, 256),
        joinRoom:   make(chan *Client),
        leaveRoom:  make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.joinRoom:
            h.mu.Lock()
            if _, exists := h.rooms[client.room]; !exists {
                h.rooms[client.room] = make(map[*Client]bool)
            }
            h.rooms[client.room][client] = true
            h.mu.Unlock()
            
            log.Printf("Client %s joined room %s", client.username, client.room)
            
        case client := <-h.leaveRoom:
            h.mu.Lock()
            if clients, exists := h.rooms[client.room]; exists {
                delete(clients, client)
                if len(clients) == 0 {
                    delete(h.rooms, client.room)
                }
            }
            h.mu.Unlock()
            
            // Notify room of user leaving
            msg := &Message{
                Type:     "user_left",
                Room:     client.room,
                Username: client.username,
            }
            h.broadcast <- msg
            
        case client := <-h.unregister:
            h.mu.Lock()
            if clients, exists := h.rooms[client.room]; exists {
                if _, ok := clients[client]; ok {
                    delete(clients, client)
                    close(client.send)
                    if len(clients) == 0 {
                        delete(h.rooms, client.room)
                    }
                }
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.RLock()
            clients := h.rooms[message.Room]
            message.UserCount = len(clients)
            h.mu.RUnlock()
            
            data, err := json.Marshal(message)
            if err != nil {
                log.Printf("Failed to marshal message: %v", err)
                continue
            }
            
            for client := range clients {
                select {
                case client.send <- data:
                default:
                    close(client.send)
                    h.mu.Lock()
                    delete(h.rooms[message.Room], client)
                    h.mu.Unlock()
                }
            }
        }
    }
}
```

### [Message Types]
```
The WebSocket chat system supports several message types:

1. **join**: Client joins a chat room
   - Client sends: {"type": "join", "room": "manga-123"}
   - Server broadcasts: {"type": "user_joined", "room": "manga-123", "username": "alice", "user_count": 5}

2. **message**: User sends chat message
   - Client sends: {"type": "message", "content": "Great chapter!"}
   - Server broadcasts: {"type": "message", "room": "manga-123", "username": "alice", "content": "Great chapter!", "timestamp": "2025-01-15T10:30:00Z"}

3. **typing**: User is typing indicator
   - Client sends: {"type": "typing", "is_typing": true}
   - Server broadcasts: {"type": "typing", "room": "manga-123", "username": "alice", "is_typing": true}

4. **user_left**: User leaves room
   - Server broadcasts: {"type": "user_left", "room": "manga-123", "username": "alice", "user_count": 4}

Messages are JSON-encoded and sent as WebSocket text frames. The server adds username and timestamp automatically, preventing spoofing.
```

### [Room Management]
```
Rooms are dynamically created when the first user joins and deleted when the last user leaves. Each room maintains its own set of connected clients. The hub's rooms map uses nested maps: outer map keys are room names, inner maps contain clients as keys with boolean values (the value is unused; the map serves as a set).

Room names follow the pattern "manga-{mangaID}" for manga-specific discussions or "general" for community chat. This enables easy room discovery and organization.

When broadcasting messages, the hub looks up all clients in the target room and sends to each. Failed sends (slow or disconnected clients) result in automatic cleanup to prevent memory leaks.
```

### [gRPC Introduction]
```
The gRPC service provides type-safe, high-performance internal APIs using Protocol Buffers for serialization. While not exposed to external clients, it demonstrates modern RPC patterns and could serve as backend-to-backend communication in a microservices architecture.
```

### [Proto Definition]
```protobuf
syntax = "proto3";

package manga;

option go_package = "mangahub/pkg/proto/manga";

// MangaService provides manga information APIs
service MangaService {
    // GetManga retrieves a single manga by ID
    rpc GetManga(GetMangaRequest) returns (MangaResponse);
    
    // SearchManga finds manga matching criteria
    rpc SearchManga(SearchRequest) returns (SearchResponse);
    
    // UpdateProgress updates user's reading progress
    rpc UpdateProgress(UpdateProgressRequest) returns (ProgressResponse);
    
    // StreamUpdates streams real-time progress updates (server streaming)
    rpc StreamUpdates(StreamRequest) returns (stream ProgressUpdate);
}

message GetMangaRequest {
    string manga_id = 1;
}

message MangaResponse {
    string id = 1;
    string title = 2;
    string author = 3;
    string artist = 4;
    repeated string genres = 5;
    string status = 6;
    int32 total_chapters = 7;
    string description = 8;
    string cover_image_url = 9;
}

message SearchRequest {
    string query = 1;
    repeated string genres = 2;
    string status = 3;
    int32 page = 4;
    int32 limit = 5;
}

message SearchResponse {
    repeated MangaResponse manga = 1;
    int32 total = 2;
}

message UpdateProgressRequest {
    string user_id = 1;
    string manga_id = 2;
    int32 current_chapter = 3;
    int32 current_page = 4;
}

message ProgressResponse {
    bool success = 1;
    string message = 2;
}

message StreamRequest {
    string user_id = 1;
}

message ProgressUpdate {
    string manga_id = 1;
    int32 current_chapter = 2;
    int32 current_page = 3;
    int64 timestamp = 4;
}
```

### [Server Implementation]
```go
package main

import (
    "context"
    "database/sql"
    "log"
    "net"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "mangahub/pkg/proto/manga"
)

type server struct {
    pb.UnimplementedMangaServiceServer
    db *sql.DB
}

func (s *server) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
    if req.MangaId == "" {
        return nil, status.Error(codes.InvalidArgument, "manga_id is required")
    }
    
    var manga pb.MangaResponse
    var genres string
    
    err := s.db.QueryRowContext(ctx,
        `SELECT id, title, author, artist, genres, status, total_chapters, description, cover_image_url
         FROM manga WHERE id = ?`,
        req.MangaId,
    ).Scan(
        &manga.Id, &manga.Title, &manga.Author, &manga.Artist,
        &genres, &manga.Status, &manga.TotalChapters,
        &manga.Description, &manga.CoverImageUrl,
    )
    
    if err == sql.ErrNoRows {
        return nil, status.Error(codes.NotFound, "manga not found")
    }
    if err != nil {
        return nil, status.Error(codes.Internal, "database error")
    }
    
    // Parse genres JSON array
    json.Unmarshal([]byte(genres), &manga.Genres)
    
    return &manga, nil
}

func (s *server) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
    // Build query with filters
    query := "SELECT id, title, author, artist, genres, status, total_chapters, description, cover_image_url FROM manga WHERE 1=1"
    args := []interface{}{}
    
    if req.Query != "" {
        query += " AND (title LIKE ? OR author LIKE ?)"
        pattern := "%" + req.Query + "%"
        args = append(args, pattern, pattern)
    }
    
    if req.Status != "" {
        query += " AND status = ?"
        args = append(args, req.Status)
    }
    
    // Add pagination
    limit := req.Limit
    if limit <= 0 || limit > 100 {
        limit = 20
    }
    offset := (req.Page - 1) * limit
    
    query += " LIMIT ? OFFSET ?"
    args = append(args, limit, offset)
    
    rows, err := s.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, status.Error(codes.Internal, "database error")
    }
    defer rows.Close()
    
    var manga []*pb.MangaResponse
    for rows.Next() {
        var m pb.MangaResponse
        var genres string
        
        err := rows.Scan(
            &m.Id, &m.Title, &m.Author, &m.Artist,
            &genres, &m.Status, &m.TotalChapters,
            &m.Description, &m.CoverImageUrl,
        )
        if err != nil {
            continue
        }
        
        json.Unmarshal([]byte(genres), &m.Genres)
        manga = append(manga, &m)
    }
    
    // Get total count
    var total int32
    countQuery := "SELECT COUNT(*) FROM manga WHERE 1=1"
    s.db.QueryRowContext(ctx, countQuery).Scan(&total)
    
    return &pb.SearchResponse{
        Manga: manga,
        Total: total,
    }, nil
}

func (s *server) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.ProgressResponse, error) {
    _, err := s.db.ExecContext(ctx,
        `INSERT INTO user_progress (user_id, manga_id, current_chapter, current_page, updated_at)
         VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
         ON CONFLICT(user_id, manga_id) DO UPDATE SET
         current_chapter = excluded.current_chapter,
         current_page = excluded.current_page,
         updated_at = CURRENT_TIMESTAMP`,
        req.UserId, req.MangaId, req.CurrentChapter, req.CurrentPage,
    )
    
    if err != nil {
        return &pb.ProgressResponse{
            Success: false,
            Message: "Failed to update progress",
        }, nil
    }
    
    return &pb.ProgressResponse{
        Success: true,
        Message: "Progress updated successfully",
    }, nil
}

func main() {
    // Initialize database
    db, err := sql.Open("sqlite3", "./mangahub.db")
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer db.Close()
    
    // Create gRPC server
    lis, err := net.Listen("tcp", ":9092")
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }
    
    grpcServer := grpc.NewServer()
    pb.RegisterMangaServiceServer(grpcServer, &server{db: db})
    
    log.Printf("gRPC server listening on :9092")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal("Failed to serve:", err)
    }
}
```

### [Client Usage]
```go
package main

import (
    "context"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    pb "mangahub/pkg/proto/manga"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:9092",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second),
    )
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer conn.Close()
    
    client := pb.NewMangaServiceClient(conn)
    
    // Example: Get manga by ID
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    manga, err := client.GetManga(ctx, &pb.GetMangaRequest{
        MangaId: "manga-123",
    })
    if err != nil {
        log.Fatal("GetManga failed:", err)
    }
    
    log.Printf("Manga: %s by %s", manga.Title, manga.Author)
    
    // Example: Search manga
    search, err := client.SearchManga(ctx, &pb.SearchRequest{
        Query: "naruto",
        Page:  1,
        Limit: 10,
    })
    if err != nil {
        log.Fatal("SearchManga failed:", err)
    }
    
    log.Printf("Found %d manga", search.Total)
    for _, m := range search.Manga {
        log.Printf("  - %s", m.Title)
    }
}
```

### [Frontend Introduction]
```
The web frontend is built with Next.js 14, providing a modern React-based user interface with server-side rendering capabilities. TypeScript ensures type safety, while Tailwind CSS enables rapid UI development with utility-first styling. The application integrates with all backend services: HTTP REST API for data operations, TCP for real-time sync, WebSocket for chat, and UDP for notifications.
```

### [Project Structure]
```
apps/web/
├── app/
│   ├── page.tsx                 # Home page (manga catalog)
│   ├── layout.tsx               # Root layout
│   ├── manga/
│   │   └── [id]/
│   │       └── page.tsx         # Manga detail page
│   ├── library/
│   │   └── page.tsx             # User library
│   ├── chat/
│   │   └── [room]/
│   │       └── page.tsx         # Chat room
│   └── auth/
│       ├── login/
│       │   └── page.tsx         # Login page
│       └── register/
│           └── page.tsx         # Registration page
├── components/
│   ├── MangaCard.tsx            # Manga display card
│   ├── MangaList.tsx            # Manga grid/list
│   ├── Navbar.tsx               # Navigation bar
│   ├── SearchBar.tsx            # Search input
│   ├── ProgressTracker.tsx      # Reading progress UI
│   └── ChatBox.tsx              # Chat interface
├── lib/
│   ├── api.ts                   # HTTP API client
│   ├── websocket.ts             # WebSocket client
│   ├── auth.ts                  # Authentication utilities
│   └── types.ts                 # TypeScript types
├── styles/
│   └── globals.css              # Global styles
└── public/
    └── images/                  # Static images
```

### [API Integration]
```typescript
// lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export interface Manga {
  id: string;
  title: string;
  author: string;
  artist: string;
  genres: string[];
  status: 'ongoing' | 'completed' | 'hiatus';
  total_chapters: number;
  description: string;
  cover_image_url: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

class MangaAPI {
  private getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };
    
    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    return headers;
  }

  async list(page: number = 1, limit: number = 20, search?: string): Promise<PaginatedResponse<Manga>> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    
    if (search) {
      params.append('search', search);
    }
    
    const response = await fetch(`${API_URL}/manga?${params}`, {
      headers: this.getHeaders(),
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch manga');
    }
    
    return response.json();
  }

  async get(id: string): Promise<Manga> {
    const response = await fetch(`${API_URL}/manga/${id}`, {
      headers: this.getHeaders(),
    });
    
    if (!response.ok) {
      throw new Error('Manga not found');
    }
    
    const data = await response.json();
    return data.data;
  }

  async getLibrary(): Promise<Manga[]> {
    const response = await fetch(`${API_URL}/library`, {
      headers: this.getHeaders(),
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch library');
    }
    
    const data = await response.json();
    return data.data;
  }

  async addToLibrary(mangaId: string): Promise<void> {
    const response = await fetch(`${API_URL}/library/${mangaId}`, {
      method: 'POST',
      headers: this.getHeaders(),
    });
    
    if (!response.ok) {
      throw new Error('Failed to add to library');
    }
  }

  async updateProgress(mangaId: string, chapter: number, page: number): Promise<void> {
    const response = await fetch(`${API_URL}/progress/${mangaId}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify({ current_chapter: chapter, current_page: page }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to update progress');
    }
  }
}

export const mangaAPI = new MangaAPI();
```

### [WebSocket Integration]
```typescript
// lib/websocket.ts
export interface ChatMessage {
  type: 'message' | 'user_joined' | 'user_left' | 'typing';
  room: string;
  username: string;
  content?: string;
  timestamp?: string;
  user_count?: number;
}

export class ChatClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  
  constructor(
    private room: string,
    private username: string,
    private onMessage: (message: ChatMessage) => void,
    private onConnectionChange: (connected: boolean) => void
  ) {}

  connect(): void {
    const wsUrl = `ws://localhost:8080/ws/chat?username=${encodeURIComponent(this.username)}`;
    
    this.ws = new WebSocket(wsUrl);
    
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;
      this.onConnectionChange(true);
      
      // Join room
      this.send({
        type: 'join',
        room: this.room,
        username: this.username,
      });
    };
    
    this.ws.onmessage = (event) => {
      try {
        const message: ChatMessage = JSON.parse(event.data);
        this.onMessage(message);
      } catch (error) {
        console.error('Failed to parse message:', error);
      }
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    this.ws.onclose = () => {
      console.log('WebSocket closed');
      this.onConnectionChange(false);
      this.attemptReconnect();
    };
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }
    
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    
    console.log(`Reconnecting in ${delay}ms... (attempt ${this.reconnectAttempts})`);
    
    setTimeout(() => {
      this.connect();
    }, delay);
  }

  sendMessage(content: string): void {
    this.send({
      type: 'message',
      room: this.room,
      username: this.username,
      content,
    });
  }

  sendTypingIndicator(isTyping: boolean): void {
    this.send({
      type: 'typing',
      room: this.room,
      username: this.username,
      is_typing: isTyping,
    });
  }

  private send(data: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    }
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
```

### [State Management]
```typescript
// Example: Chat component with state management
'use client';

import { useState, useEffect, useRef } from 'react';
import { ChatClient, ChatMessage } from '@/lib/websocket';

export default function ChatRoom({ room }: { room: string }) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [connected, setConnected] = useState(false);
  const [userCount, setUserCount] = useState(0);
  const chatClient = useRef<ChatClient | null>(null);
  
  useEffect(() => {
    const username = localStorage.getItem('username') || 'Anonymous';
    
    chatClient.current = new ChatClient(
      room,
      username,
      (message) => {
        setMessages(prev => [...prev, message]);
        if (message.user_count) {
          setUserCount(message.user_count);
        }
      },
      (isConnected) => {
        setConnected(isConnected);
      }
    );
    
    chatClient.current.connect();
    
    return () => {
      chatClient.current?.disconnect();
    };
  }, [room]);

  const handleSend = () => {
    if (inputValue.trim() && chatClient.current) {
      chatClient.current.sendMessage(inputValue);
      setInputValue('');
    }
  };

  return (
    <div className="flex flex-col h-screen">
      <div className="bg-blue-600 text-white p-4">
        <h1 className="text-xl font-bold">Chat Room: {room}</h1>
        <p className="text-sm">{connected ? `${userCount} users online` : 'Disconnected'}</p>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {messages.map((msg, index) => (
          <div key={index} className={`p-2 rounded ${
            msg.type === 'message' ? 'bg-gray-100' : 'bg-yellow-50 italic'
          }`}>
            {msg.type === 'message' ? (
              <>
                <span className="font-bold">{msg.username}:</span> {msg.content}
              </>
            ) : (
              <span>{msg.username} {msg.type === 'user_joined' ? 'joined' : 'left'}</span>
            )}
          </div>
        ))}
      </div>
      
      <div className="p-4 border-t">
        <div className="flex gap-2">
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSend()}
            placeholder="Type a message..."
            className="flex-1 px-4 py-2 border rounded"
            disabled={!connected}
          />
          <button
            onClick={handleSend}
            disabled={!connected}
            className="px-6 py-2 bg-blue-600 text-white rounded disabled:bg-gray-400"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}
```

---

I'll continue with Chapter 5, 6, References, and Appendices in the next message. Shall I proceed?
