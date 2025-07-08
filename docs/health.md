# Health Check API

This document covers the health check endpoint in the Chirpy API.

## Overview

The Health Check API provides a simple endpoint to verify that the Chirpy server is running and responding to requests. This is commonly used for monitoring, load balancers, and deployment verification.

## Base URL

The health check endpoint is prefixed with `/api`

## Endpoints

### GET /api/healthz

Check if the server is running and healthy.

**Authentication:** Not required

**Response (200 OK):**

```
OK
```

**Response Headers:**

```
Content-Type: text/plain; charset=utf-8
```

**Example:**

```bash
curl http://localhost:8080/api/healthz
```

**Response:**

```
OK
```

## Usage

### Monitoring and Alerts

Use this endpoint to monitor server health:

```bash
# Basic health check
curl -f http://localhost:8080/api/healthz

# With timeout for monitoring scripts
curl -f --max-time 5 http://localhost:8080/api/healthz

# Check status code only
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/healthz
```

### Load Balancer Configuration

Configure your load balancer to use this endpoint:

```
Health Check Path: /api/healthz
Expected Response: 200 OK
Expected Content: "OK"
```

### Docker Health Checks

Add to your Dockerfile:

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:8080/api/healthz || exit 1
```

### Kubernetes Probes

Configure liveness and readiness probes:

```yaml
livenessProbe:
  httpGet:
    path: /api/healthz
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /api/healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Response Details

### Success Response

- **Status Code:** 200 OK
- **Content-Type:** text/plain; charset=utf-8
- **Body:** `OK`

### What This Endpoint Checks

The health check endpoint verifies:

- Server process is running
- HTTP server is accepting connections
- Basic request handling is functional

### What This Endpoint Does NOT Check

- Database connectivity
- External service availability
- Application-specific functionality
- Memory/CPU usage
- Disk space

## Error Scenarios

### Server Not Running

```bash
curl http://localhost:8080/api/healthz
# curl: (7) Failed to connect to localhost port 8080: Connection refused
```

### Server Overloaded

The endpoint may respond slowly or timeout if the server is under heavy load.

## Best Practices

### Monitoring Setup

1. **Regular Checks:** Monitor every 30-60 seconds
2. **Timeout:** Set reasonable timeout (5-10 seconds)
3. **Retries:** Allow 2-3 retries before marking as unhealthy
4. **Alerting:** Alert on consecutive failures

### Integration Examples

#### Shell Script Monitoring

```bash
#!/bin/bash
if curl -f --max-time 5 http://localhost:8080/api/healthz > /dev/null 2>&1; then
    echo "Server is healthy"
    exit 0
else
    echo "Server is unhealthy"
    exit 1
fi
```

#### Python Health Check

```python
import requests
import sys

try:
    response = requests.get('http://localhost:8080/api/healthz', timeout=5)
    if response.status_code == 200 and response.text.strip() == 'OK':
        print("Server is healthy")
        sys.exit(0)
    else:
        print("Server returned unexpected response")
        sys.exit(1)
except requests.exceptions.RequestException as e:
    print(f"Health check failed: {e}")
    sys.exit(1)
```

## Related Endpoints

For more comprehensive monitoring, consider also checking:

- [Admin Metrics](admin.md#get-adminmetrics) - For system usage statistics
- Database connectivity (not exposed via API)
- Key application endpoints with authentication

## Troubleshooting

### Common Issues

1. **Connection Refused:**

   - Server process not running
   - Port 8080 not available
   - Firewall blocking connections

2. **Timeout:**

   - Server overloaded
   - Network connectivity issues
   - Process deadlock

3. **Wrong Response:**
   - Server running but application code failing
   - Proxy or load balancer interference

### Debugging Steps

1. **Check server logs**
2. **Verify port availability:** `netstat -tlnp | grep 8080`
3. **Test locally:** `curl localhost:8080/api/healthz`
4. **Check system resources:** CPU, memory, disk usage
