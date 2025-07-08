# Admin API

This document covers all administrative endpoints in the Chirpy API.

## Overview

The Admin API provides system management capabilities including metrics monitoring and data reset functionality. These endpoints are typically used for system administration and monitoring.

## Base URL

All admin endpoints are prefixed with `/admin`

## Endpoints

### GET /admin/metrics

Get system metrics including file server hit count.

**Authentication:** Not required

**Response (200 OK):**

```html
<html>
  <body>
    <h1>Admin</h1>
    <p>Chirpy has been visited [HIT_COUNT] times!</p>
  </body>
</html>
```

**Error Responses:**

- `500 Internal Server Error` - Server error

**Example:**

```bash
curl http://localhost:8080/admin/metrics
```

**Response Example:**

```html
<html>
  <body>
    <h1>Admin</h1>
    <p>Chirpy has been visited 1247 times!</p>
  </body>
</html>
```

### POST /admin/reset

Reset system data and metrics.

**Authentication:** Not required

**Platform Restriction:** Only available when `PLATFORM` environment variable is set to `"dev"`

**Request Body:** Empty

**Response (200 OK):**

```json
{
  "message": "Hits reset to 0 and file server hits reset to 0"
}
```

**Error Responses:**

- `401 Unauthorized` - Not available in production environment
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/admin/reset
```

## Metrics System

### File Server Hits

The system tracks the number of times static files are served through the file server:

- Counter is incremented for each request to `/app/*` endpoints
- Displayed in the admin metrics page
- Reset to 0 when using the reset endpoint

### Metrics Display

- Metrics are displayed as an HTML page
- Shows total file server hits since last reset
- Designed for easy monitoring and debugging

## Development vs Production

### Development Mode

When `PLATFORM=dev`:

- All admin endpoints are available
- Reset functionality works normally
- Useful for development and testing

### Production Mode

When `PLATFORM` is not set to `"dev"`:

- Metrics endpoint remains available
- Reset endpoint returns `401 Unauthorized`
- Prevents accidental data loss in production

## Security Considerations

### No Authentication Required

- Admin endpoints currently don't require authentication
- In production, consider adding authentication/authorization
- Metrics endpoint exposes system usage information

### Platform Protection

- Reset endpoint is protected by environment variable
- Only works in development mode
- Prevents accidental data loss

## Usage Examples

### Monitoring System Health

```bash
# Check current metrics
curl http://localhost:8080/admin/metrics

# Example response shows system activity
# Use this to monitor application usage
```

### Development Workflow

```bash
# During development/testing
curl -X POST http://localhost:8080/admin/reset

# Reset all metrics and data
# Useful for clean testing environments
```

### Integration with Monitoring Tools

```bash
# You can parse the HTML response to extract metrics
curl -s http://localhost:8080/admin/metrics | grep -o '[0-9]\+' | head -1
# Returns just the hit count number
```

## Error Handling

Admin endpoints return appropriate HTTP status codes:

- `200 OK` - Successful operation
- `401 Unauthorized` - Operation not allowed (production reset)
- `500 Internal Server Error` - Server error

## Future Enhancements

Potential improvements for the admin system:

- JSON response format for metrics
- Authentication/authorization
- More detailed system metrics
- Database statistics
- User count and activity metrics
- API endpoint usage statistics

## File Server Integration

The admin system integrates with the file server middleware:

- Tracks requests to `/app/*` endpoints
- Increments counter for each static file served
- Provides visibility into frontend usage

## Environment Configuration

Required environment variables:

- `PLATFORM` - Set to `"dev"` to enable reset functionality

Optional for production:

- Consider adding authentication keys
- Add monitoring system integration
