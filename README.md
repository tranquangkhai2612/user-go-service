# User Go Service

A RESTful API service for user management built with Go.

## Project Structure

```
user-go-service/
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   └── service/         # Business logic
├── pkg/                 # Public packages
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
go mod download
```

### Running the Service

```bash
go run cmd/server/main.go
```

### Building

```bash
go build -o bin/user-service cmd/server/main.go
```

## API Endpoints

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## License

MIT
