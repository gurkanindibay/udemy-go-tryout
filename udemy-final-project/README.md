# Udemy Go Final Project - Event Management API (REST + gRPC)

A comprehensive event management system built with Go that provides both RESTful API and gRPC services for managing events, user authentication, and event registrations.

## Overview

This project is a complete event management system that allows users to:
- Register and authenticate
- Create, read, update, and delete events
- Register for events
- View their event registrations

The API uses SQLite as the database, JWT for authentication, and follows RESTful conventions.

## Features

- **User Management**: User registration and login with JWT authentication
- **Event Management**: Full CRUD operations for events
- **Event Registration**: Users can register for events and view their registrations
- **Authentication**: JWT-based authentication for protected routes
- **Database**: SQLite database with proper schema and relationships
- **Dual API Support**: Both RESTful HTTP API and gRPC services
- **RESTful API**: Clean REST endpoints following standard conventions

## Technology Stack

- **Language**: Go 1.25.0
- **Web Framework**: Gin (for REST API)
- **RPC Framework**: gRPC (for gRPC services)
- **Database**: SQLite (modernc.org/sqlite driver)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Protocol Buffers**: For gRPC service definitions
- **Dependency Injection**: samber/do (for service management)
- **JSON**: Standard library with Gin bindings

## Dependency Injection

This project implements Dependency Injection (DI) using the `samber/do` library to manage service dependencies and improve code maintainability, testability, and decoupling.

### Why Dependency Injection?

- **Testability**: Services can be easily mocked for unit testing
- **Maintainability**: Changes to service implementations don't affect dependent code
- **Decoupling**: Components are loosely coupled and follow the Dependency Inversion Principle
- **Flexibility**: Easy to swap implementations or add new features

### How DI is Implemented

#### 1. Service Interfaces
All business logic is abstracted behind interfaces defined in `services/interfaces.go`:
- `UserService` - User management operations
- `EventService` - Event management operations
- `AuthService` - Authentication operations

#### 2. Service Implementations
Concrete implementations are provided in `services/implementations.go`:
- `NewUserService()` - Creates user service instance
- `NewEventService()` - Creates event service instance
- `NewAuthService()` - Creates auth service instance

#### 3. DI Container
The DI container is set up in `di/container.go`:
```go
func NewContainer() *Container {
    injector := do.New()

    // Register services with named values
    do.ProvideNamedValue(injector, "userService", services.NewUserService())
    do.ProvideNamedValue(injector, "eventService", services.NewEventService())
    do.ProvideNamedValue(injector, "authService", services.NewAuthService())

    return &Container{Injector: injector}
}
```

#### 4. Service Usage
Services are injected into route handlers and gRPC servers:

**REST Routes** (`routes/routes.go`):
```go
func SetupRoutes(router *gin.Engine, container *di.Container) {
    // Inject services into route handlers
    userService := container.GetUserService()
    eventService := container.GetEventService()
    authService := container.GetAuthService()
    
    // Services are available as package-level variables
    // for use in individual route files
}
```

**gRPC Servers** (`grpc/auth/server.go`, `grpc/event/server.go`):
```go
func NewAuthServer(container *di.Container) *AuthServer {
    return &AuthServer{
        authService: container.GetAuthService(),
    }
}
```

### Benefits in This Project

1. **Clean Architecture**: Business logic is separated from HTTP/gRPC concerns
2. **Easy Testing**: Services can be mocked for comprehensive unit tests
3. **Service Sharing**: Package-level variables allow services to be shared across route files
4. **Future Extensibility**: New services can be added without modifying existing code
5. **Configuration Management**: Services can be configured differently for different environments

### DI Container Methods

The container provides getter methods for each service:
- `GetUserService()` - Returns the user service instance
- `GetEventService()` - Returns the event service instance
- `GetAuthService()` - Returns the auth service instance

This ensures type safety and centralized service management throughout the application.

## Prerequisites

- Go 1.25.0 or later
- Git (for cloning the repository)

## Installation & Setup

1. **Clone the repository** (if not already done):
   ```bash
   git clone https://github.com/gurkanindibay/udemy-go-tryout.git
   cd udemy-go-tryout/udemy-final-project
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## Database

The application uses SQLite with the following schema:

### Tables

- **users**: Stores user information
  - `id` (INTEGER, PRIMARY KEY)
  - `email` (TEXT, NOT NULL)
  - `password` (TEXT, NOT NULL, hashed)

- **events**: Stores event information
  - `id` (INTEGER, PRIMARY KEY)
  - `name` (TEXT, NOT NULL)
  - `description` (TEXT, NOT NULL)
  - `location` (TEXT, NOT NULL)
  - `date_time` (TEXT, NOT NULL)
  - `user_id` (INTEGER, FOREIGN KEY to users.id)

- **registrations**: Links users to events they've registered for
  - `id` (INTEGER, PRIMARY KEY)
  - `event_id` (INTEGER, FOREIGN KEY to events.id)
  - `user_id` (INTEGER, FOREIGN KEY to users.id)

The database file `events.db` is created automatically when the application starts.

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user

#### Events
- `GET /events` - Get all events
- `GET /events/:id` - Get event by ID

### Protected Endpoints (Authentication Required)

#### Events
- `POST /events` - Create a new event
- `PUT /events/:id` - Update an event
- `DELETE /events/:id` - Delete an event

#### Registrations
- `POST /events/:id/register` - Register for an event
- `DELETE /events/:id/register` - Cancel event registration
- `GET /users/:id/registrations` - Get user's event registrations

## gRPC Services

The application also provides gRPC services running on port `50051`. The gRPC services mirror the functionality of the REST API but use Protocol Buffers for efficient communication.

### gRPC Services Available

#### AuthService
- `Register(RegisterRequest) returns (RegisterResponse)` - Register a new user
- `Login(LoginRequest) returns (LoginResponse)` - Authenticate user and return JWT token

#### EventService
- `GetEvents(GetEventsRequest) returns (GetEventsResponse)` - Get all events
- `GetEvent(GetEventRequest) returns (GetEventResponse)` - Get event by ID
- `CreateEvent(CreateEventRequest) returns (CreateEventResponse)` - Create a new event
- `UpdateEvent(UpdateEventRequest) returns (UpdateEventResponse)` - Update an existing event
- `DeleteEvent(DeleteEventRequest) returns (DeleteEventResponse)` - Delete an event
- `RegisterForEvent(RegisterForEventRequest) returns (RegisterForEventResponse)` - Register for an event
- `CancelRegistration(CancelRegistrationRequest) returns (CancelRegistrationResponse)` - Cancel event registration
- `GetUserRegistrations(GetUserRegistrationsRequest) returns (GetUserRegistrationsResponse)` - Get user's registrations

### gRPC Client Example

A sample gRPC client is provided in `client/grpc_client.go`. To run it:

```bash
go run client/grpc_client.go
```

### Protocol Buffer Definitions

The gRPC services are defined in:
- `proto/auth.proto` - Authentication service
- `proto/event.proto` - Event management service

Generated Go code is in:
- `proto/auth/` - Auth service generated code
- `proto/event/` - Event service generated code

## Request/Response Examples

### User Registration
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### User Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

### Create Event
```http
POST /events
Content-Type: application/json
Authorization: Bearer <your-jwt-token>

{
  "name": "Sample Event",
  "description": "This is a sample event",
  "location": "Sample Location",
  "date_time": "2023-10-10T10:00:00Z"
}
```

### Get All Events
```http
GET /events
```

Response:
```json
[
  {
    "id": 1,
    "name": "Sample Event",
    "description": "This is a sample event",
    "location": "Sample Location",
    "date_time": "2023-10-10T10:00:00Z",
    "user_id": 1
  }
]
```

## Testing

### Using HTTP Test Files

The project includes pre-configured `.http` files in the `api-test/` directory for testing all endpoints. You can use these with VS Code's REST Client extension or similar tools.

Available test files:
- `auth-register.http` - User registration
- `login.http` - User login
- `create-event.http` - Create new event
- `get-events.http` - Get all events
- `get-events-by-id.http` - Get event by ID
- `update-event.http` - Update event
- `delete-event.http` - Delete event
- `register.http` - Register for event
- `delete-register.http` - Cancel registration

### Manual Testing with curl

1. **Register a user**:
   ```bash
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Login**:
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

3. **Create an event** (replace TOKEN with actual JWT):
   ```bash
   curl -X POST http://localhost:8080/events \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer TOKEN" \
     -d '{
       "name": "Test Event",
       "description": "A test event",
       "location": "Test Location",
       "date_time": "2023-10-10T10:00:00Z"
     }'
   ```

### Testing Tips

1. Always register/login first to get a JWT token
2. Include the JWT token in the Authorization header for protected routes
3. Use tools like Postman, Insomnia, or VS Code REST Client for easier testing
4. Check the database file `events.db` to verify data persistence

## API Documentation

The API includes interactive Swagger UI documentation that allows you to:

- View all available endpoints
- See request/response schemas
- Test API endpoints directly from the browser
- Download the OpenAPI specification

### Accessing Swagger UI

1. Start the server: `go run main.go`
2. Open your browser and navigate to: `http://localhost:8080/swagger/index.html`
3. You'll see the complete API documentation with interactive testing capabilities

### Features

- **Interactive Testing**: Test all endpoints directly from the browser
- **Request/Response Examples**: See sample requests and responses
- **Authentication**: JWT token input for protected endpoints
- **Schema Validation**: View data models and validation rules

## Project Structure

```
udemy-final-project/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── events.db              # SQLite database (created automatically)
├── api-test/              # HTTP test files
│   ├── auth-register.http
│   ├── create-event.http
│   ├── delete-event.http
│   ├── delete-register.http
│   ├── get-events-by-id.http
│   ├── get-events-by-id.http
│   ├── get-events.http
│   ├── login.http
│   ├── register.http
│   └── update-event.http
├── client/
│   └── grpc_client.go     # Sample gRPC client
├── db/
│   └── db.go              # Database initialization and connection
├── grpc/
│   ├── auth/
│   │   └── server.go      # gRPC auth service implementation
│   └── event/
│       └── server.go      # gRPC event service implementation
├── middlewares/
│   └── auth.go            # JWT authentication middleware
├── models/
│   ├── event.go           # Event model and database operations
│   └── user.go            # User model and authentication
├── proto/
│   ├── auth.proto         # Auth service protobuf definition
│   ├── event.proto        # Event service protobuf definition
│   ├── auth/              # Generated auth protobuf code
│   └── event/             # Generated event protobuf code
├── routes/
│   ├── events.go          # Event-related REST routes
│   ├── registers.go       # Registration-related REST routes
│   ├── routes.go          # Main REST route setup
│   └── users.go           # User-related REST routes
└── utils/
    ├── hash.go            # Password hashing utilities
    └── jwt.go             # JWT token utilities
```

## Security Features

- **Password Hashing**: Uses bcrypt for secure password storage
- **JWT Authentication**: Stateless authentication with expiration
- **Input Validation**: Gin binding validation for request data
- **SQL Injection Protection**: Prepared statements for all database queries

## Development

### Building for Production

```bash
go build -o event-api main.go
```

### Running Tests

Currently, the project doesn't have unit tests implemented, but you can add them using Go's testing framework.

### Environment Variables

The application currently uses hardcoded values but can be extended to use environment variables for:
- Database path
- JWT secret key
- Server port
- etc.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is for educational purposes as part of the Udemy Go course.