# Chirpy API Documentation

Welcome to the comprehensive API documentation for Chirpy - a Twitter-like social media platform built with Go.

## Quick Start

1. **Health Check**: Start by testing if the server is running

   ```bash
   curl http://localhost:8080/api/healthz
   ```

2. **Create User**: Register a new account

   ```bash
   curl -X POST http://localhost:8080/api/users \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "password123"}'
   ```

3. **Login**: Get authentication tokens

   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "password123"}'
   ```

4. **Create Chirp**: Post your first message
   ```bash
   curl -X POST http://localhost:8080/api/chirps \
     -H "Authorization: Bearer <your_token>" \
     -H "Content-Type: application/json" \
     -d '{"body": "Hello, Chirpy world!"}'
   ```

## API Documentation

### Core APIs

#### [Authentication API](auth.md)

- User login and logout
- JWT token management
- Refresh token handling
- **Key Endpoints:**
  - `POST /api/login` - User authentication
  - `POST /api/refresh` - Refresh access tokens
  - `POST /api/revoke` - Revoke refresh tokens

#### [Users API](users.md)

- User registration and management
- Profile updates
- Premium status tracking
- **Key Endpoints:**
  - `POST /api/users` - Create new user
  - `PUT /api/users` - Update user profile

#### [Chirps API](chirps.md)

- Create, read, and delete chirps
- Content filtering and validation
- Author filtering and sorting
- **Key Endpoints:**
  - `POST /api/chirps` - Create new chirp
  - `GET /api/chirps` - List all chirps
  - `GET /api/chirps/{id}` - Get specific chirp
  - `DELETE /api/chirps/{id}` - Delete chirp

### System APIs

#### [Admin API](admin.md)

- System metrics and monitoring
- Development utilities
- **Key Endpoints:**
  - `GET /admin/metrics` - View system metrics
  - `POST /admin/reset` - Reset system data (dev only)

#### [Webhooks API](webhooks.md)

- External service integrations
- Premium subscription management
- **Key Endpoints:**
  - `POST /api/polka/webhooks` - Polka payment webhooks

#### [Health Check API](health.md)

- Server health monitoring
- Load balancer integration
- **Key Endpoints:**
  - `GET /api/healthz` - Health check endpoint

## Authentication Flow

```
1. Register User    → POST /api/users
2. Login           → POST /api/login (returns access + refresh tokens)
3. Use Access Token → Include in Authorization header
4. Refresh Token   → POST /api/refresh (when access token expires)
5. Logout          → POST /api/revoke (invalidate refresh token)
```

## Common HTTP Status Codes

- **200 OK** - Successful GET/PUT requests
- **201 Created** - Successful POST requests
- **204 No Content** - Successful DELETE requests
- **400 Bad Request** - Invalid request data
- **401 Unauthorized** - Authentication required/failed
- **403 Forbidden** - Access denied
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

## Base URLs

- **API Endpoints**: `http://localhost:8080/api/`
- **Admin Endpoints**: `http://localhost:8080/admin/`
- **Static Files**: `http://localhost:8080/app/`

## Request/Response Format

### Request Headers

```
Content-Type: application/json
Authorization: Bearer <token>  # For authenticated endpoints
```

### Error Response Format

Most endpoints return plain text error messages:

```
HTTP/1.1 400 Bad Request
Content-Type: text/plain

Invalid request data
```

### Success Response Format

JSON responses for data endpoints:

```json
{
  "id": "uuid",
  "field": "value",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production use.

## CORS

Cross-Origin Resource Sharing (CORS) is not configured. Add CORS middleware if needed for web applications.

## Environment Variables

Required for API functionality:

```bash
DB_URL=postgres://user:password@localhost:5432/chirpy
JWT_SECRET=your-secret-key
POLKA_KEY=your-polka-webhook-key
PLATFORM=dev  # For admin reset functionality
```

## Testing the API

### Using curl

All examples in the documentation use curl. Make sure you have curl installed:

```bash
curl --version
```

### Using Postman

Import the following base configuration:

- Base URL: `http://localhost:8080`
- Add Authorization header for authenticated requests
- Set Content-Type to `application/json` for POST/PUT requests

### Using HTTPie

Alternative to curl:

```bash
# Install HTTPie
pip install httpie

# Example usage
http POST localhost:8080/api/users email=user@example.com password=password123
```

## API Client Libraries

Consider implementing client libraries for:

- JavaScript/TypeScript
- Python
- Go
- Mobile SDKs (iOS/Android)

## Support

For questions about the API:

1. Check the specific endpoint documentation
2. Verify your request format matches the examples
3. Check server logs for error details
4. Ensure all required environment variables are set

## Contributing

When adding new endpoints:

1. Update the corresponding documentation file
2. Include request/response examples
3. Document all error cases
4. Add authentication requirements
5. Update this index file if adding new API sections
