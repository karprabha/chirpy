# Chirps API

This document covers all chirp-related endpoints in the Chirpy API.

## Overview

The Chirps API allows users to create, read, and delete short messages called "chirps" (similar to tweets). All chirp operations require authentication except for reading chirps.

## Base URL

All chirp endpoints are prefixed with `/api`

## Endpoints

### POST /api/chirps

Create a new chirp.

**Authentication:** Required (Bearer token)

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "body": "This is my first chirp! 🐦"
}
```

**Response (201 Created):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "body": "This is my first chirp! 🐦",
  "user_id": "987fcdeb-51a2-43d7-b456-426614174000",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON, missing body, or chirp too long (>140 characters)
- `401 Unauthorized` - Invalid, expired, or missing access token
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello, Chirpy world!"}'
```

### GET /api/chirps

Get all chirps, optionally filtered by author.

**Authentication:** Not required

**Query Parameters:**

- `author_id` (optional) - UUID of the user to filter chirps by
- `sort` (optional) - Sort order: `asc` (default) or `desc`

**Response (200 OK):**

```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "body": "This is my first chirp! 🐦",
    "user_id": "987fcdeb-51a2-43d7-b456-426614174000",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "id": "234e5678-e89b-12d3-a456-426614174001",
    "body": "Another chirp from the same user",
    "user_id": "987fcdeb-51a2-43d7-b456-426614174000",
    "created_at": "2023-01-01T01:00:00Z",
    "updated_at": "2023-01-01T01:00:00Z"
  }
]
```

**Error Responses:**

- `400 Bad Request` - Invalid author_id format
- `500 Internal Server Error` - Server error

**Examples:**

```bash
# Get all chirps (ascending order by default)
curl http://localhost:8080/api/chirps

# Get chirps by specific author
curl "http://localhost:8080/api/chirps?author_id=987fcdeb-51a2-43d7-b456-426614174000"

# Get chirps in descending order (newest first)
curl "http://localhost:8080/api/chirps?sort=desc"

# Get chirps by author in descending order
curl "http://localhost:8080/api/chirps?author_id=987fcdeb-51a2-43d7-b456-426614174000&sort=desc"
```

### GET /api/chirps/{id}

Get a specific chirp by ID.

**Authentication:** Not required

**Path Parameters:**

- `id` (required) - UUID of the chirp to retrieve

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "body": "This is my first chirp! 🐦",
  "user_id": "987fcdeb-51a2-43d7-b456-426614174000",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid ID format
- `404 Not Found` - Chirp not found
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl http://localhost:8080/api/chirps/123e4567-e89b-12d3-a456-426614174000
```

### DELETE /api/chirps/{id}

Delete a specific chirp.

**Authentication:** Required (Bearer token)

**Authorization:** Users can only delete their own chirps

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Path Parameters:**

- `id` (required) - UUID of the chirp to delete

**Response (204 No Content):**
Empty response body

**Error Responses:**

- `400 Bad Request` - Invalid ID format
- `401 Unauthorized` - Invalid, expired, or missing access token
- `403 Forbidden` - User doesn't own the chirp
- `404 Not Found` - Chirp not found
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X DELETE http://localhost:8080/api/chirps/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Chirp Model

### Chirp Object

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "body": "This is my first chirp! 🐦",
  "user_id": "987fcdeb-51a2-43d7-b456-426614174000",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Fields:**

- `id` (UUID) - Unique identifier for the chirp
- `body` (string) - The chirp content (max 140 characters)
- `user_id` (UUID) - ID of the user who created the chirp
- `created_at` (timestamp) - When the chirp was created
- `updated_at` (timestamp) - When the chirp was last updated

## Content Rules

### Length Limit

- Chirps must be **140 characters or less**
- Exceeding this limit returns a `400 Bad Request` error

### Profanity Filter

The API automatically filters profanity from chirp content:

- Filtered words: "kerfuffle", "sharbert", "fornax"
- Filtered words are replaced with "\*\*\*\*"
- Case-insensitive filtering

**Example:**

```json
// Input
{"body": "What a kerfuffle this is!"}

// Output
{"body": "What a **** this is!"}
```

## Sorting and Filtering

### Sort Options

- `asc` (default) - Oldest chirps first
- `desc` - Newest chirps first

### Author Filtering

- Use `author_id` query parameter to get chirps from a specific user
- Must be a valid UUID format
- Returns empty array if author has no chirps

## Security and Authorization

### Authentication

- **Creating chirps**: Requires valid access token
- **Deleting chirps**: Requires valid access token
- **Reading chirps**: No authentication required

### Authorization

- Users can only delete their own chirps
- Attempting to delete another user's chirp returns `403 Forbidden`

## Error Handling

Chirp endpoints return appropriate HTTP status codes:

- `200 OK` - Successful retrieval
- `201 Created` - Successful creation
- `204 No Content` - Successful deletion
- `400 Bad Request` - Invalid request data or validation errors
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - User doesn't own the resource
- `404 Not Found` - Chirp not found
- `500 Internal Server Error` - Server error

## Example Workflows

### Creating and Managing Chirps

1. **Create a chirp:**

   ```bash
   curl -X POST http://localhost:8080/api/chirps \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{"body": "My first chirp!"}'
   ```

2. **Get all chirps:**

   ```bash
   curl http://localhost:8080/api/chirps
   ```

3. **Get your chirps:**

   ```bash
   curl "http://localhost:8080/api/chirps?author_id=<your_user_id>"
   ```

4. **Delete a chirp:**
   ```bash
   curl -X DELETE http://localhost:8080/api/chirps/<chirp_id> \
     -H "Authorization: Bearer <access_token>"
   ```

### Building a Timeline

```bash
# Get latest chirps for a timeline view
curl "http://localhost:8080/api/chirps?sort=desc"

# Get chirps from a specific user
curl "http://localhost:8080/api/chirps?author_id=<user_id>&sort=desc"
```
