# JWT Authentication Setup Guide

## Overview

Your user service now includes JWT authentication with the following features:

- User registration and login
- Protected API endpoints
- Token-based authentication
- Swagger UI for testing

## Database Migration

First, run the migration to add the password field to your database:

```bash
mysql -u root -p < scripts/add_password_field.sql
```

## Configuration

Update your `.env` file with a secure JWT secret:

```bash
JWT_SECRET=your-very-secure-secret-key-change-this
```

## Running the Application

```bash
go run cmd/server/main.go
```

The server will start on http://localhost:8080

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Register a New User

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securePassword123"
}
```

#### Login

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "...",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### Protected Endpoints (Require Authentication)

All user management endpoints now require a JWT token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

#### Get All Users

```bash
GET /api/v1/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Get User by ID

```bash
GET /api/v1/users/{id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Create User

```bash
POST /api/v1/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "email": "newuser@example.com",
  "name": "Jane Smith",
  "password": "password123"
}
```

#### Update User

```bash
PUT /api/v1/users/{id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "email": "updated@example.com",
  "name": "Updated Name"
}
```

#### Delete User

```bash
DELETE /api/v1/users/{id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Using Swagger UI

1. Open http://localhost:8080/swagger/
2. Click on `/api/v1/auth/register` or `/api/v1/auth/login`
3. Click "Try it out"
4. Enter your credentials and execute
5. Copy the token from the response
6. Click the "Authorize" button at the top of the page
7. Enter `Bearer <your-token>` (include the word "Bearer" followed by a space)
8. Click "Authorize"
9. Now you can test all protected endpoints!

## Token Details

- **Expiration**: Tokens expire after 24 hours
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Claims**: User ID and email are included in the token

## Security Notes

1. **Change the JWT_SECRET** in production to a strong, random string
2. Passwords are hashed using bcrypt with default cost (10)
3. Always use HTTPS in production to protect tokens in transit
4. Tokens are not stored server-side (stateless authentication)
5. Consider implementing token refresh mechanism for production

## Testing the Flow

1. **Register a new user:**

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","name":"Test User","password":"test123"}'
   ```

2. **Login and get token:**

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"test123"}'
   ```

3. **Use token to access protected endpoint:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/users \
     -H "Authorization: Bearer YOUR_TOKEN_HERE"
   ```

## Error Responses

- `401 Unauthorized`: Missing or invalid token
- `400 Bad Request`: Invalid request body
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

## JWT Authentication Flow

```
Client                          Server
  |                               |
  |  1. POST /auth/register       |
  |------------------------------>|
  |  {email, name, password}      |
  |                               |
  |  2. Return JWT + User         |
  |<------------------------------|
  |  {token, user}                |
  |                               |
  |  3. GET /users                |
  |  Authorization: Bearer token  |
  |------------------------------>|
  |                               |
  |  4. Validate token            |
  |     Extract user info         |
  |                               |
  |  5. Return protected data     |
  |<------------------------------|
```
