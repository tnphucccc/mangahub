# MangaHub API Documentation

Complete REST API documentation for the MangaHub HTTP server.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Authentication Endpoints](#authentication-endpoints)
- [Manga Endpoints](#manga-endpoints)
- [User Endpoints](#user-endpoints-protected)
- [Health Check](#health-check)
- [Error Handling](#error-handling)
- [Examples](#complete-api-testing-example)

---

## Overview

### Base URL

```
http://localhost:8080/api/v1
```

### Content Type

All requests and responses use JSON format:

```
Content-Type: application/json
```

### Authentication

Protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens are valid for 7 days by default and are obtained through the login or register endpoints.

---

## Authentication Endpoints

### Register a New User

Create a new user account and receive a JWT token.

**Endpoint:**

```http
POST /api/v1/auth/register
```

**Request Body:**

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Validation Rules:**

- `username`: required, 3-32 characters
- `email`: required, valid email format
- `password`: required, minimum 6 characters

**Success Response (201 Created):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-11-27T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwidXNlcm5hbWUiOiJqb2huZG9lIiwiZXhwIjoxNzM3OTcyNjAwLCJuYmYiOjE3Mzc5NjU0MDAsImlhdCI6MTczNzk2NTQwMH0.abcdef123456..."
}
```

**Error Responses:**

`400 Bad Request` - Invalid request body or validation failed:

```json
{
  "error": "Invalid request body"
}
```

`409 Conflict` - Username or email already exists:

```json
{
  "error": "username or email already exists"
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "alice123"
  }'
```

---

### Login

Authenticate an existing user and receive a JWT token.

**Endpoint:**

```http
POST /api/v1/auth/login
```

**Request Body:**

```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Success Response (200 OK):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-11-27T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**

`400 Bad Request` - Invalid request body:

```json
{
  "error": "Invalid request body"
}
```

`401 Unauthorized` - Invalid credentials:

```json
{
  "error": "Invalid username or password"
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "alice123"
  }'
```

---

## Manga Endpoints

### Search Manga

Search for manga by title, author, genre, or status with pagination support.

**Endpoint:**

```http
GET /api/v1/manga
```

**Query Parameters:**

| Parameter | Type    | Required | Description                                                      | Default |
| --------- | ------- | -------- | ---------------------------------------------------------------- | ------- |
| `title`   | string  | No       | Filter by manga title (case-insensitive, partial match)          | -       |
| `author`  | string  | No       | Filter by author name (case-insensitive, partial match)          | -       |
| `genre`   | string  | No       | Filter by genre                                                  | -       |
| `status`  | string  | No       | Filter by status (`ongoing`, `completed`, `hiatus`, `cancelled`) | -       |
| `limit`   | integer | No       | Number of results (max: 100)                                     | 20      |
| `offset`  | integer | No       | Pagination offset                                                | 0       |

**Success Response (200 OK):**

```json
{
  "manga": [
    {
      "id": "manga-002",
      "title": "Naruto",
      "author": "Masashi Kishimoto",
      "genres": ["Action", "Adventure", "Martial Arts"],
      "status": "completed",
      "total_chapters": 700,
      "description": "The story follows Naruto Uzumaki, a young ninja who seeks recognition from his peers and dreams of becoming the Hokage.",
      "cover_image_url": "https://example.com/naruto.jpg",
      "created_at": "2025-11-27T03:08:20Z",
      "updated_at": "2025-11-27T03:08:20Z"
    }
  ],
  "count": 1
}
```

**Examples:**

```bash
# Search by title
curl "http://localhost:8080/api/v1/manga?title=naruto"

# Search by author
curl "http://localhost:8080/api/v1/manga?author=kishimoto"

# Filter by genre
curl "http://localhost:8080/api/v1/manga?genre=Action"

# Filter by status
curl "http://localhost:8080/api/v1/manga?status=ongoing"

# Combined filters with pagination
curl "http://localhost:8080/api/v1/manga?title=naruto&genre=Action&status=ongoing&limit=10&offset=0"
```

---

### Get All Manga

Retrieve all manga with pagination.

**Endpoint:**

```http
GET /api/v1/manga/all
```

**Query Parameters:**

| Parameter | Type    | Required | Description                  | Default |
| --------- | ------- | -------- | ---------------------------- | ------- |
| `limit`   | integer | No       | Number of results (max: 100) | 20      |
| `offset`  | integer | No       | Pagination offset            | 0       |

**Success Response (200 OK):**

```json
{
  "manga": [
    {
      "id": "manga-001",
      "title": "One Piece",
      "author": "Eiichiro Oda",
      "genres": ["Action", "Adventure", "Fantasy"],
      "status": "ongoing",
      "total_chapters": 1100,
      "description": "The story follows Monkey D. Luffy, a young man whose body gained the properties of rubber after unintentionally eating a Devil Fruit.",
      "cover_image_url": "https://example.com/onepiece.jpg",
      "created_at": "2025-11-27T03:08:20Z",
      "updated_at": "2025-11-27T03:08:20Z"
    }
  ],
  "count": 10
}
```

**Example:**

```bash
curl "http://localhost:8080/api/v1/manga/all?limit=20&offset=0"
```

---

### Get Manga by ID

Retrieve detailed information about a specific manga.

**Endpoint:**

```http
GET /api/v1/manga/:id
```

**Path Parameters:**

- `id` (required): Manga ID

**Success Response (200 OK):**

```json
{
  "manga": {
    "id": "manga-001",
    "title": "One Piece",
    "author": "Eiichiro Oda",
    "genres": ["Action", "Adventure", "Fantasy"],
    "status": "ongoing",
    "total_chapters": 1100,
    "description": "The story follows Monkey D. Luffy, a young man whose body gained the properties of rubber after unintentionally eating a Devil Fruit.",
    "cover_image_url": "https://example.com/onepiece.jpg",
    "created_at": "2025-11-27T03:08:20Z",
    "updated_at": "2025-11-27T03:08:20Z"
  }
}
```

**Error Responses:**

`404 Not Found` - Manga does not exist:

```json
{
  "error": "Manga not found"
}
```

**Example:**

```bash
curl "http://localhost:8080/api/v1/manga/manga-001"
```

---

## User Endpoints (Protected)

All user endpoints require authentication via JWT token.

### Get Current User Profile

Retrieve the profile of the currently authenticated user.

**Endpoint:**

```http
GET /api/v1/users/me
```

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-11-27T10:30:00Z"
  }
}
```

**Error Responses:**

`401 Unauthorized` - Missing or invalid token:

```json
{
  "error": "Authorization header required"
}
```

or

```json
{
  "error": "Invalid or expired token"
}
```

**Example:**

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  "http://localhost:8080/api/v1/users/me"
```

---

### Get User's Manga Library

Retrieve all manga in the user's library with reading progress.

**Endpoint:**

```http
GET /api/v1/users/library
```

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "library": [
    {
      "user_id": "user-testuser",
      "manga_id": "manga-001",
      "current_chapter": 50,
      "status": "reading",
      "rating": {
        "Int64": 8,
        "Valid": true
      },
      "started_at": {
        "Time": "2025-11-20T10:00:00Z",
        "Valid": true
      },
      "completed_at": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "updated_at": "2025-11-27T03:08:20Z",
      "manga": {
        "id": "manga-001",
        "title": "One Piece",
        "author": "Eiichiro Oda",
        "genres": ["Action", "Adventure", "Fantasy"],
        "status": "ongoing",
        "total_chapters": 1100,
        "description": "The story follows Monkey D. Luffy...",
        "cover_image_url": "https://example.com/onepiece.jpg",
        "created_at": "2025-11-27T03:08:20Z",
        "updated_at": "2025-11-27T03:08:20Z"
      }
    }
  ],
  "count": 1
}
```

**Reading Status Values:**

- `reading`: Currently reading
- `completed`: Finished reading
- `plan_to_read`: Planning to read
- `on_hold`: Temporarily stopped
- `dropped`: Abandoned

**Error Responses:**

`401 Unauthorized` - Missing or invalid token

**Example:**

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  "http://localhost:8080/api/v1/users/library"
```

---

### Add Manga to Library

Add a manga to the user's library with initial status and progress.

**Endpoint:**

```http
POST /api/v1/users/library
```

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "manga_id": "manga-003",
  "status": "plan_to_read",
  "current_chapter": 0
}
```

**Request Fields:**

| Field             | Type    | Required | Description                         |
| ----------------- | ------- | -------- | ----------------------------------- |
| `manga_id`        | string  | Yes      | ID of the manga to add              |
| `status`          | string  | Yes      | Reading status (see options below)  |
| `current_chapter` | integer | No       | Current chapter number (default: 0) |

**Status Options:**

- `reading`: Currently reading
- `completed`: Finished reading
- `plan_to_read`: Planning to read
- `on_hold`: Temporarily stopped
- `dropped`: Abandoned

**Success Response (201 Created):**

```json
{
  "message": "Manga added to library"
}
```

**Error Responses:**

`400 Bad Request` - Invalid request body:

```json
{
  "error": "Invalid request body"
}
```

`401 Unauthorized` - Missing or invalid token

`404 Not Found` - Manga does not exist:

```json
{
  "error": "Manga not found"
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/users/library \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "manga_id": "manga-003",
    "status": "reading",
    "current_chapter": 0
  }'
```

---

### Get Reading Progress

Retrieve reading progress for a specific manga.

**Endpoint:**

```http
GET /api/v1/users/progress/:manga_id
```

**Headers:**

```
Authorization: Bearer <token>
```

**Path Parameters:**

- `manga_id` (required): Manga ID

**Success Response (200 OK):**

```json
{
  "progress": {
    "user_id": "user-testuser",
    "manga_id": "manga-001",
    "current_chapter": 50,
    "status": "reading",
    "rating": {
      "Int64": 8,
      "Valid": true
    },
    "started_at": {
      "Time": "2025-11-20T10:00:00Z",
      "Valid": true
    },
    "completed_at": {
      "Time": "0001-01-01T00:00:00Z",
      "Valid": false
    },
    "updated_at": "2025-11-27T03:08:20Z"
  }
}
```

**Error Responses:**

`401 Unauthorized` - Missing or invalid token

`404 Not Found` - Progress not found (manga not in library):

```json
{
  "error": "Progress not found"
}
```

**Example:**

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  "http://localhost:8080/api/v1/users/progress/manga-001"
```

---

### Update Reading Progress

Update reading progress for a manga in the user's library.

**Endpoint:**

```http
PUT /api/v1/users/progress/:manga_id
```

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters:**

- `manga_id` (required): Manga ID

**Request Body:**

```json
{
  "current_chapter": 100,
  "status": "reading",
  "rating": 9
}
```

**Request Fields:**

| Field             | Type    | Required | Description            | Constraints         |
| ----------------- | ------- | -------- | ---------------------- | ------------------- |
| `current_chapter` | integer | Yes      | Current chapter number | 0 to total_chapters |
| `status`          | string  | No       | Reading status         | See status options  |
| `rating`          | integer | No       | User rating            | 1-10                |

**Success Response (200 OK):**

```json
{
  "message": "Progress updated"
}
```

**Error Responses:**

`400 Bad Request` - Invalid request body:

```json
{
  "error": "Invalid request body"
}
```

or invalid chapter number:

```json
{
  "error": "Invalid chapter number"
}
```

`401 Unauthorized` - Missing or invalid token

`404 Not Found` - Manga not in library:

```json
{
  "error": "Manga not in library"
}
```

**Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/users/progress/manga-001 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "current_chapter": 100,
    "rating": 8
  }'
```

---

## Health Check

### Check API Health

Check if the API server is running and database is connected.

**Endpoint:**

```http
GET /health
```

**Success Response (200 OK):**

```json
{
  "status": "ok",
  "service": "MangaHub HTTP API Server",
  "version": "1.0.0"
}
```

**Error Response (500 Internal Server Error):**

```json
{
  "status": "unhealthy",
  "service": "MangaHub HTTP API Server",
  "error": "database connection failed"
}
```

**Example:**

```bash
curl "http://localhost:8080/health"
```

---

## Error Handling

### Standard Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### HTTP Status Codes

| Status Code | Meaning               | Description                                        |
| ----------- | --------------------- | -------------------------------------------------- |
| 200         | OK                    | Request succeeded                                  |
| 201         | Created               | Resource created successfully                      |
| 400         | Bad Request           | Invalid request format or validation failed        |
| 401         | Unauthorized          | Missing or invalid authentication token            |
| 404         | Not Found             | Resource not found                                 |
| 409         | Conflict              | Resource already exists (e.g., duplicate username) |
| 500         | Internal Server Error | Server-side error                                  |

### Common Error Scenarios

**Missing Authorization Header:**

```json
{
  "error": "Authorization header required"
}
```

**Invalid Token Format:**

```json
{
  "error": "Invalid authorization header format"
}
```

**Expired or Invalid Token:**

```json
{
  "error": "Invalid or expired token"
}
```

**Validation Error:**

```json
{
  "error": "Invalid request body"
}
```

---

## Complete API Testing Example

Here's a complete workflow for testing all API endpoints:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

# 1. Register a new user
echo "1. Registering new user..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "alice123"
  }')

echo $REGISTER_RESPONSE | jq .

# Extract token from response
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# 2. Login with existing user
echo -e "\n2. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "alice123"
  }')

echo $LOGIN_RESPONSE | jq .

# Update token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

# 3. Get current user profile
echo -e "\n3. Getting user profile..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/users/me" | jq .

# 4. Search for manga
echo -e "\n4. Searching for manga..."
curl -s "$BASE_URL/manga?title=one%20piece" | jq .

# 5. Get manga by ID
echo -e "\n5. Getting manga by ID..."
curl -s "$BASE_URL/manga/manga-001" | jq .

# 6. Add manga to library
echo -e "\n6. Adding manga to library..."
curl -s -X POST $BASE_URL/users/library \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "manga_id": "manga-001",
    "status": "reading",
    "current_chapter": 0
  }' | jq .

# 7. Update reading progress
echo -e "\n7. Updating reading progress..."
curl -s -X PUT $BASE_URL/users/progress/manga-001 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_chapter": 50,
    "rating": 9
  }' | jq .

# 8. Get reading progress
echo -e "\n8. Getting reading progress..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/users/progress/manga-001" | jq .

# 9. Get user's library
echo -e "\n9. Getting user library..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/users/library" | jq .

# 10. Health check
echo -e "\n10. Checking API health..."
curl -s "http://localhost:8080/health" | jq .

echo -e "\n✅ All API endpoints tested successfully!"
```

Save this as `test-api.sh`, make it executable with `chmod +x test-api.sh`, and run it to test all endpoints.

---

## Rate Limiting & Performance

Currently, there are no rate limits implemented. For production use, consider:

- Implementing rate limiting middleware
- Adding request logging
- Setting up monitoring and metrics
- Configuring connection pooling appropriately

## Security Considerations

- Always use HTTPS in production
- Store JWT secrets securely (environment variables, not in code)
- Implement rate limiting to prevent brute force attacks
- Add input sanitization for all user inputs
- Consider implementing refresh tokens for long-lived sessions
- Add request logging for security auditing

---

**Last Updated**: 2025-11-27
**API Version**: 1.0.0
**Protocol**: HTTP REST API (25 points)
