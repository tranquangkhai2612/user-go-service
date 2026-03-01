# Role-Based Access Control (RBAC)

This application now supports role-based access control with two roles: `admin` and `user`.

## Roles

- **`user`** (default): Can view users (GET endpoints)
- **`admin`**: Full access to all endpoints (GET, POST, PUT, DELETE)

## Database Migration

Run the migration to add the `role` column:

```bash
# If using Docker MySQL
docker exec -i mysql-container mysql -u root -p<password> user_service < scripts/add_role_field.sql

# If using local MySQL
mysql -u root -p < scripts/add_role_field.sql
```

## Registration

### Register as Regular User (default)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "Regular User",
    "password": "password123"
  }'
```

Response includes user with `"role": "user"` by default.

### Register as Admin

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "name": "Admin User",
    "password": "admin123",
    "role": "admin"
  }'
```

Response includes user with `"role": "admin"`.

## JWT Token

The JWT token includes the role in its claims:

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "user",
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890
}
```

## Protected Endpoints

### Available to Both Roles (user + admin)

- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get user by ID

### Admin Only

- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/{id}` - Update a user
- `DELETE /api/v1/users/{id}` - Delete a user

## Testing Role Permissions

Use the provided test script:

```bash
# Make sure the server is running
./scripts/test_roles.sh
```

This script will:

1. Register an admin user
2. Register a regular user
3. Test that both can GET users
4. Test that only admin can DELETE users
5. Test that only admin can CREATE users

## Example Usage

### 1. Register and login as admin

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","name":"Admin","password":"admin123","role":"admin"}'

# Login (or use token from registration)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'
```

Save the returned token.

### 2. Delete a user (admin only)

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

**Expected**: `200 OK` with success message

### 3. Try to delete as regular user

```bash
# First register/login as regular user
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"user123"}'

# Try to delete
curl -X DELETE http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer <USER_TOKEN>"
```

**Expected**: `403 Forbidden` with error message:

```json
{
  "error": "Access denied: insufficient permissions"
}
```

## Error Responses

### 401 Unauthorized

- Missing or invalid JWT token
- Token expired

```json
{
  "error": "Invalid or expired token"
}
```

### 403 Forbidden

- Valid token but insufficient permissions (wrong role)

```json
{
  "error": "Access denied: insufficient permissions"
}
```

## Implementation Details

### Models

- `User` struct includes `Role` field (string: "user" or "admin")
- Role constants: `models.RoleUser` and `models.RoleAdmin`

### JWT Claims

- `JWTClaims` struct includes `Role` field
- Role is embedded in the token during authentication

### Middleware

- `Authenticate`: Validates JWT token (required for all protected routes)
- `RequireRole(roles...)`: Checks if user has one of the specified roles

### Route Protection

In `cmd/server/main.go`:

```go
// All protected routes require authentication
protected := api.PathPrefix("").Subrouter()
protected.Use(jwtMiddleware.Authenticate)

// Admin-only routes
adminRoutes := protected.PathPrefix("").Subrouter()
adminRoutes.Use(jwtMiddleware.RequireRole("admin"))
adminRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
adminRoutes.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
adminRoutes.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

// Routes available to all authenticated users
protected.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
```

## Security Considerations

1. **Role Assignment**: Currently, users can specify their role during registration. In production, you should:
   - Remove role from `RegisterRequest`
   - Always default new users to `"user"` role
   - Only allow admins to promote users to admin via a separate endpoint

2. **Admin Promotion**: Add an admin-only endpoint to promote users:

   ```go
   PUT /api/v1/users/{id}/promote
   ```

3. **Role Validation**: The system validates roles are either "user" or "admin" during registration.
