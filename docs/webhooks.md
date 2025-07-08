# Webhooks API

This document covers webhook endpoints in the Chirpy API.

## Overview

The Webhooks API allows external services to send notifications to Chirpy about events like user upgrades. Currently, it supports Polka payment service webhooks for premium subscription upgrades.

## Base URL

All webhook endpoints are prefixed with `/api`

## Endpoints

### POST /api/polka/webhooks

Receive webhooks from Polka payment service for user upgrades.

**Authentication:** API Key required

**Request Headers:**

```
Authorization: ApiKey <polka_api_key>
Content-Type: application/json
```

**Request Body:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

**Response (204 No Content):**
Empty response body

**Error Responses:**

- `400 Bad Request` - Invalid JSON format
- `401 Unauthorized` - Invalid or missing API key
- `404 Not Found` - User not found
- `500 Internal Server Error` - Server error

**Example:**

```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Authorization: ApiKey your-polka-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "event": "user.upgraded",
    "data": {
      "user_id": "123e4567-e89b-12d3-a456-426614174000"
    }
  }'
```

## Webhook Events

### user.upgraded

Notifies Chirpy when a user upgrades to Chirpy Red (premium subscription).

**Event Data:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

**Processing:**

1. Validates the API key
2. Checks if the user exists
3. Updates user's `is_chirpy_red` field to `true`
4. Returns success response

**Effect:**

- User gains premium features
- `is_chirpy_red` field becomes `true` in user profile
- User can access premium functionality

### Other Events

Currently, only `user.upgraded` events are processed. Other event types will:

- Return `204 No Content` (success)
- No processing occurs
- Allows for future event type expansion

## Authentication

### API Key Authentication

Webhooks use API key authentication:

- Set via `POLKA_KEY` environment variable
- Must be included in `Authorization` header
- Format: `Authorization: ApiKey <your_api_key>`

**Example Header:**

```
Authorization: ApiKey pk_test_abc123def456
```

### Security

- API key must match the configured `POLKA_KEY`
- Invalid keys return `401 Unauthorized`
- Missing keys return `401 Unauthorized`

## Error Handling

### Webhook-Specific Errors

**400 Bad Request:**

- Malformed JSON payload
- Missing required fields
- Invalid user ID format

**401 Unauthorized:**

- Missing Authorization header
- Invalid API key
- Malformed API key format

**404 Not Found:**

- User ID doesn't exist in system
- User was deleted after payment processed

**500 Internal Server Error:**

- Database connection issues
- Failed to update user record

## Integration Guide

### Setting up Polka Webhooks

1. **Configure Environment:**

   ```bash
   export POLKA_KEY="your-polka-api-key"
   ```

2. **Register Webhook URL:**
   Register `https://your-domain.com/api/polka/webhooks` with Polka

3. **Test Webhook:**
   ```bash
   curl -X POST https://your-domain.com/api/polka/webhooks \
     -H "Authorization: ApiKey your-polka-api-key" \
     -H "Content-Type: application/json" \
     -d '{
       "event": "user.upgraded",
       "data": {
         "user_id": "valid-user-uuid"
       }
     }'
   ```

### Webhook Verification

To verify webhooks are working:

1. **Check User Status Before:**

   ```bash
   curl http://localhost:8080/api/users/123e4567-e89b-12d3-a456-426614174000
   # Should show "is_chirpy_red": false
   ```

2. **Send Webhook:**

   ```bash
   curl -X POST http://localhost:8080/api/polka/webhooks \
     -H "Authorization: ApiKey your-key" \
     -H "Content-Type: application/json" \
     -d '{"event": "user.upgraded", "data": {"user_id": "123e4567-e89b-12d3-a456-426614174000"}}'
   ```

3. **Check User Status After:**
   ```bash
   curl http://localhost:8080/api/users/123e4567-e89b-12d3-a456-426614174000
   # Should show "is_chirpy_red": true
   ```

## Webhook Payload Examples

### Valid Upgrade Event

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### Ignored Event (Future Use)

```json
{
  "event": "user.downgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

## Best Practices

### Webhook Handling

1. **Idempotency:** Webhooks should be idempotent

   - Multiple calls with same data should have same effect
   - Upgrading an already upgraded user is safe

2. **Retry Logic:** Implement retry on client side

   - Handle temporary failures gracefully
   - Use exponential backoff for retries

3. **Validation:** Always validate incoming data
   - Check user exists before processing
   - Validate UUID format

### Security

1. **API Key Management:**

   - Keep API keys secure and confidential
   - Rotate keys regularly
   - Use environment variables, not hardcoded values

2. **HTTPS Only:**

   - Always use HTTPS in production
   - Protect webhook data in transit

3. **Rate Limiting:**
   - Consider implementing rate limiting
   - Prevent abuse of webhook endpoints

## Future Enhancements

Potential webhook improvements:

- Support for user downgrades
- Webhook signature validation
- Webhook event logging
- Batch webhook processing
- Webhook retry mechanisms
- More granular premium feature controls

## Environment Configuration

Required environment variables:

```bash
POLKA_KEY=your-polka-api-key-here
```

The API key should be provided by Polka when setting up webhook integration.
