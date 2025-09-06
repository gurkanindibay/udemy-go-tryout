# Udemy Go Final Project - Event Management REST API

A RESTful API built with Go and Gin framework for managing events, user authentication, and event registrations.

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
- **RESTful API**: Clean REST endpoints following standard conventions

## Technology Stack

- **Language**: Go 1.25.0
- **Web Framework**: Gin
- **Database**: SQLite (modernc.org/sqlite driver)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **JSON**: Standard library with Gin bindings

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
│   ├── get-events.http
│   ├── login.http
│   ├── register.http
│   └── update-event.http
├── db/
│   └── db.go              # Database initialization and connection
├── middlewares/
│   └── auth.go            # JWT authentication middleware
├── models/
│   ├── event.go           # Event model and database operations
│   └── user.go            # User model and authentication
├── routes/
│   ├── events.go          # Event-related routes
│   ├── registers.go       # Registration-related routes
│   ├── routes.go          # Main route setup
│   └── users.go           # User-related routes
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