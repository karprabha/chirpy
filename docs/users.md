# Users API

This document covers all user-related endpoints in the Chirpy API.

## Overview

The Users API allows you to create and manage user accounts. All user operations except registration require authentication.

## Base URL

All user endpoints are prefixed with `/api`

## Endpoints

### POST /api/users

Create a new user account.

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON, missing email/password, or validation errors
- `500 Internal Server Error` - Server error (possibly duplicate email)

**Example:**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "newuser@example.com", "password": "securepassword123"}'
```

### PUT /api/users

Update an existing user's information.

**Authentication:** Required (Bearer token)

**Request Headers:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "newemail@example.com",
  "is_chirpy_red": false,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON, missing email/password, or validation errors
- `401 Unauthorized` - Invalid, expired, or missing access token
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"email": "updated@example.com", "password": "newpassword123"}'
```

## User Model

### User Object

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

**Fields:**

- `id` (UUID) - Unique identifier for the user
- `email` (string) - User's email address (must be unique)
- `is_chirpy_red` (boolean) - Premium subscription status
- `created_at` (timestamp) - When the user account was created
- `updated_at` (timestamp) - When the user account was last updated

## Validation Rules

### Email

- Must be a valid email format
- Must be unique across all users
- Required for both registration and updates

### Password

- No minimum length enforced by API (implement client-side validation)
- Required for both registration and updates
- Stored as bcrypt hash (never returned in responses)

## Premium Features (Chirpy Red)

The `is_chirpy_red` field indicates whether a user has premium features:

- `false` - Standard user
- `true` - Premium user with Chirpy Red features

**Note:** Users cannot directly update their premium status through the API. Premium upgrades are handled through the [Webhooks API](webhooks.md).

## Security Considerations

- **Password Security**: Passwords are hashed using bcrypt and never returned in API responses
- **Email Uniqueness**: The system enforces unique email addresses
- **Authentication**: User updates require a valid access token
- **User Isolation**: Users can only update their own information

## Error Handling

User endpoints return appropriate HTTP status codes:

- `200 OK` - Successful update
- `201 Created` - Successful registration
- `400 Bad Request` - Invalid request data or validation errors
- `401 Unauthorized` - Authentication required or failed
- `500 Internal Server Error` - Server error

Error responses include a plain text message in the response body.

## Integration with Authentication

After creating a user account:

1. Use the [Authentication API](auth.md#post-apilogin) to log in
2. Use the returned access token for authenticated requests
3. Update user information using the access token

## Example User Flow

1. **Registration:**

   ```bash
   curl -X POST http://localhost:8080/api/users \
     -H "Content-Type: application/json" \
     -d '{"email": "john@example.com", "password": "password123"}'
   ```

2. **Login:**

   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email": "john@example.com", "password": "password123"}'
   ```

3. **Update Profile:**
   ```bash
   curl -X PUT http://localhost:8080/api/users \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{"email": "john.doe@example.com", "password": "newpassword123"}'
   ```
