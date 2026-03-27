
# API Reference

Base URL: `http://localhost:8080`

## Public Routes

### GET /health
Check if the server is running.

**Response**
```json
{ "status": "healthy", 
  "database": "connected" }
```

---

## Auth Routes

No authentication required.

### POST /api/verify-token
Exchange a Firebase ID token for a custom JWT access token + refresh token. Creates the user in the database if they don't exist yet.

**Request**
```json
{
  "token": "<firebase_id_token>"
}
```

**Response**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 604800,
  "user": {
    "id": 1,
    "firebase_uid": "abc123",
    "email": "user@example.com",
    "firstname": "Tran",
    "lastname": "Tran"
  }
}
```

---

### POST /api/refresh-token
Get a new access token using a refresh token.

**Request**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 604800
}
```

---

## Protected Routes

Requires `Authorization: Bearer <access_token>` header.

---

### POST /api/fcm/register
Register a device FCM token for push notifications.

**Request**
```json
{
  "token": "<fcm_device_token>"
}
```

---

### DELETE /api/fcm/delete
Remove a device FCM token.

**Request**
```json
{
  "token": "<fcm_device_token>"
}
```

---

### POST /api/fcm/test
Send a test push notification to the authenticated user's devices.

---

## Admin Routes

Requires `Authorization: Bearer <access_token>` header. User must have `is_admin = true` in the database.

Returns `403 Forbidden` if the user is not an admin.

### GET /api/admin/users
Get all users.

**Response**
```json
{
  "data": [
    {
      "id": 1,
      "firstname": "Tran",
      "lastname": "Tran",
      "email": "user@example.com",
      "is_admin": true
    }
  ]
}
```

---

### PATCH /api/admin/users/promote
Promote a user to admin by email.

**Request**
```json
{
  "email": "user@example.com"
}
```

**Response**
```json
{
  "data": {
    "id": 1,
    "firstname": "Tran",
    "lastname": "Tran",
    "email": "user@example.com",
    "is_admin": true
  }
}
```

---

### PATCH /api/admin/users/demote
Demote a user from admin by email.

**Request**
```json
{
  "email": "user@example.com"
}
```

**Response**
```json
{
  "data": {
    "id": 1,
    "firstname": "Tran",
    "lastname": "Tran",
    "email": "user@example.com",
    "is_admin": false
  }
}
```

---

## Error Responses

| Status | Meaning |
|--------|---------|
| 400 | Bad request — missing or invalid input |
| 401 | Unauthorized — missing or invalid token |
| 403 | Forbidden — not an admin |
| 404 | Not found |
| 500 | Internal server error |