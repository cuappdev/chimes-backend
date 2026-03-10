# chimes-backend

Backend for the Chimes app. Built with Go, Gin, GORM, PostgreSQL, and Firebase.

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [Docker](https://www.docker.com/) (for running PostgreSQL locally)
- Firebase project with a service account key

## Local Setup

1. Clone the repo
2. Copy the Firebase service account key to the root as `service-account-key.json`
3. Create a `.env` file in the root (see Environment Variables below)
4. Start the database:
   ```bash
   docker compose up db -d
   ```
5. Run the app:
   ```bash
   go run main.go
   ```

The server starts on `http://localhost:8080`.

## Environment Variables

Create a `.env` file in the root:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<your_password>
DB_NAME=chimes
DB_SSLMODE=disable
JWT_SECRET=<any_random_string>
IMAGE_TAG=latest
```

## Auth Flow

Authentication uses Firebase + custom JWTs:

1. Client signs in with Firebase → receives a Firebase ID token 
2. Client calls `POST /api/verify-token` with the Firebase token → receives an access token + refresh token (expire in 7 days)
3. Client uses the access token for all subsequent API calls via `Authorization: Bearer <access_token>`
4. When the access token expires, call `POST /api/refresh-token` to get a new one
5. When the refresh token expires, repeat from step 1

## Promoting a User to Admin

Admin status must be set manually in the database. There is no API endpoint for this.

1. Find the user's ID:
   ```sql
   SELECT id, email FROM users;
   ```
2. Set them as admin:
   ```sql
   UPDATE users SET is_admin = true WHERE id = <user_id>;
   ```

Admin users can access all routes under `/api/admin/`.

## Testing with Postman (No Frontend)

Since there is no frontend yet, use the Firebase REST API to get a Firebase ID token:

**Step 1 — Get Firebase ID token**
```
POST https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=<FIREBASE_WEB_API_KEY>

{
  "email": "user@example.com",
  "password": "yourpassword",
  "returnSecureToken": true
}
```
Copy the `idToken` from the response.

**Step 2 — Exchange for JWT**
```
POST http://localhost:8080/api/verify-token

{
  "token": "<idToken from step 1>"
}
```
Copy the full `access_token` from the response.

**Step 3 — Call API**
```
GET http://localhost:8080/api/admin/users
Authorization: Bearer <access_token from step 2>
```

The Firebase Web API Key is in Firebase Console → Project Settings → General.

## API Endpoints

See [API.md](API.md) for full endpoint documentation.
