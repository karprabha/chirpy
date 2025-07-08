# Authentication API

This document covers all authentication-related endpoints in the Chirpy API.

## Overview

Chirpy uses JWT (JSON Web Tokens) for authentication with a refresh token system:

- **Access tokens** expire in 1 hour and are used for API requests
- **Refresh tokens** expire in 60 days and are used to obtain new access tokens
- Passwords are hashed using bcrypt

## Base URL

All authentication endpoints are prefixed with `/api`

## Endpoints

### POST /api/login

Authenticate a user and receive access and refresh tokens.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "abc123def456..."
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON or missing email/password
- `401 Unauthorized` - Incorrect email or password
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### POST /api/refresh

Get a new access token using a refresh token.

**Request Headers:**

```
Authorization: Bearer <refresh_token>
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**

- `401 Unauthorized` - Invalid, expired, or missing refresh token
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer abc123def456..."
```

### POST /api/revoke

Revoke (invalidate) a refresh token.

**Request Headers:**

```
Authorization: Bearer <refresh_token>
```

**Response (204 No Content):**
Empty response body

**Error Responses:**

- `401 Unauthorized` - Invalid or missing refresh token
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer abc123def456..."
```

## Authentication Flow

### 1. User Registration

First, create a user account using the [Users API](users.md#post-apiusers).

### 2. Login

Use the `/api/login` endpoint with email and password to receive:

- Access token (for API requests)
- Refresh token (for renewing access tokens)

### 3. Making Authenticated Requests

Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### 4. Token Refresh

When the access token expires (after 1 hour), use the refresh token to get a new access token via `/api/refresh`.

### 5. Logout

Revoke the refresh token using `/api/revoke` to log out the user.

## Security Considerations

- **Password Requirements**: While not enforced by the API, use strong passwords
- **Token Storage**: Store tokens securely (avoid localStorage for sensitive apps)
- **HTTPS**: Always use HTTPS in production
- **Token Expiration**: Access tokens expire in 1 hour, refresh tokens in 60 days
- **Password Hashing**: Passwords are hashed using bcrypt with default cost

## Error Handling

All authentication endpoints return appropriate HTTP status codes:

- `200 OK` - Success
- `204 No Content` - Success (no response body)
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication failed
- `500 Internal Server Error` - Server error

Error responses include a plain text message in the response body.
