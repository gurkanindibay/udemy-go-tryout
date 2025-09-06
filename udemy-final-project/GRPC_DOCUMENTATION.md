# gRPC Services Documentation

This document provides comprehensive information about the gRPC services implemented in the Udemy Go Final Project, including implementation details, usage examples, and setup instructions.

## Overview

The application now supports both REST API and gRPC services for event management and user authentication. The gRPC services provide the same functionality as the REST API but with improved performance and type safety through Protocol Buffers.

## Architecture

```
Client Applications
├── REST API (HTTP/JSON) ──┐
│   ├── Gin Web Framework  │
│   └── JSON Serialization │
└──────────────────────────┼─► Application Server
                           │
                           ├── gRPC Services ──┐
                           │   ├── Protocol Buffers
                           │   └── Binary Serialization
                           └── SQLite Database
```

## Service Definitions

### Protocol Buffer Files

#### 1. Authentication Service (`proto/auth.proto`)

See: `proto/auth.proto`

#### 2. Event Management Service (`proto/event.proto`)

See: `proto/event.proto`

## Server Implementation

### Server Structure

```
grpc/
├── auth/
│   └── server.go          # AuthService implementation
└── event/
    └── server.go          # EventService implementation
```

### Authentication Service Implementation

See: `grpc/auth/server.go`

### Event Service Implementation

See: `grpc/event/server.go`

### Server Setup

See: `main.go` - `startGRPCServer()` function

## Client Usage

### Go Client Example

See: `client/grpc_client.go`

This file contains a complete example of:
- Connecting to the gRPC server
- Using the AuthService for user registration and login
- Using the EventService for CRUD operations
- Proper error handling and response processing

## Setup and Installation

### Prerequisites

1. **Go 1.25.0 or later**
2. **Protocol Buffers Compiler (protoc)**
3. **Go plugins for protoc**

### Installation Steps

1. **Install Protocol Buffers Compiler:**
   ```bash
   # Download and install protoc
   curl -L https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-win64.zip -o protoc.zip
   unzip protoc.zip
   # Add protoc.exe to your PATH
   ```

2. **Install Go Plugins:**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **Generate Code from Proto Files:**
   ```bash
   # Set PATH to include Go binaries
   export PATH=$PATH:$(go env GOPATH)/bin

   # Generate code
   protoc --go_out=. --go-grpc_out=. --proto_path=. --proto_path=./include proto/*.proto
   ```

4. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

5. **Run the Server:**
   ```bash
   go run main.go
   ```

   See: `main.go` for server startup implementation

## API Reference

### Authentication Service

#### Register
**Request:** `RegisterRequest`
- `email` (string): User email address
- `password` (string): User password

**Response:** `RegisterResponse`
- `user` (User): Registered user information

#### Login
**Request:** `LoginRequest`
- `email` (string): User email address
- `password` (string): User password

**Response:** `LoginResponse`
- `token` (string): JWT authentication token
- `message` (string): Success message

### Event Service

#### GetEvents
**Request:** `GetEventsRequest` (empty)

**Response:** `GetEventsResponse`
- `events` ([]Event): List of all events

#### GetEvent
**Request:** `GetEventRequest`
- `id` (int64): Event ID

**Response:** `GetEventResponse`
- `event` (Event): Event details

#### CreateEvent
**Request:** `CreateEventRequest`
- `name` (string): Event name
- `description` (string): Event description
- `location` (string): Event location
- `date_time` (Timestamp): Event date and time

**Response:** `CreateEventResponse`
- `event` (Event): Created event

#### UpdateEvent
**Request:** `UpdateEventRequest`
- `id` (int64): Event ID
- `name` (string): Updated event name
- `description` (string): Updated event description
- `location` (string): Updated event location
- `date_time` (Timestamp): Updated event date and time

**Response:** `UpdateEventResponse`
- `event` (Event): Updated event

#### DeleteEvent
**Request:** `DeleteEventRequest`
- `id` (int64): Event ID

**Response:** `DeleteEventResponse` (empty)

#### RegisterForEvent
**Request:** `RegisterForEventRequest`
- `event_id` (int64): Event ID to register for

**Response:** `RegisterForEventResponse` (empty)

#### CancelRegistration
**Request:** `CancelRegistrationRequest`
- `event_id` (int64): Event ID to cancel registration for

**Response:** `CancelRegistrationResponse` (empty)

#### GetUserRegistrations
**Request:** `GetUserRegistrationsRequest` (empty)

**Response:** `GetUserRegistrationsResponse`
- `events` ([]Event): List of events user is registered for

## Data Types

### User
- `id` (int64): User ID
- `email` (string): User email address
- `password` (string): User password (only in requests)

### Event
- `id` (int64): Event ID
- `name` (string): Event name
- `description` (string): Event description
- `location` (string): Event location
- `date_time` (Timestamp): Event date and time
- `user_id` (int64): ID of user who created the event

## Error Handling

gRPC services return errors in the following cases:

- **Invalid credentials**: When login fails
- **Event not found**: When trying to access non-existent event
- **Permission denied**: When user tries to modify event they don't own
- **Database errors**: When database operations fail
- **Validation errors**: When request data is invalid

## Security Considerations

### Authentication
- JWT tokens are used for authentication
- Passwords are hashed using bcrypt
- User sessions are stateless

### Authorization
- Event ownership is enforced
- Users can only modify their own events
- Registration operations require authentication

### Transport Security
- Use TLS in production (`grpc.WithTransportCredentials()`)
- Never send sensitive data in plaintext
- Validate all input data

## Performance Optimization

### Connection Management
- Use connection pooling for multiple requests
- Implement connection reuse
- Handle connection failures gracefully

### Streaming
- Consider using streaming for large datasets
- Implement pagination for list operations
- Use bidirectional streaming for real-time features

### Caching
- Cache frequently accessed data
- Implement proper cache invalidation
- Use Redis or similar for distributed caching

## Testing

### Unit Tests
See: Test files in the project (when implemented)

### Integration Tests
See: `client/grpc_client.go` for integration testing examples

### HTTP Test Files
See: `api-test/` directory for REST API testing

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure server is running on correct port
   - Check firewall settings
   - Verify network connectivity

2. **Protobuf Generation Errors**
   - Ensure protoc is installed and in PATH
   - Check Go plugins are installed
   - Verify proto file syntax

3. **Import Errors**
   - Run `go mod tidy`
   - Check GOPATH and module paths
   - Ensure generated files are in correct location

### Debug Mode
Enable gRPC reflection for debugging:
```go
reflection.Register(grpcServer)
```

Use tools like `grpcui` for interactive testing:
```bash
grpcui -plaintext localhost:50051
```

## Migration from REST to gRPC

### Benefits of gRPC
- **Type Safety**: Protocol Buffers ensure type safety
- **Performance**: Binary serialization is faster than JSON
- **Streaming**: Support for streaming operations
- **Code Generation**: Auto-generated client/server code
- **Multi-language**: Support for multiple programming languages

### Migration Strategy
1. Keep REST API for backward compatibility
2. Add gRPC services alongside REST
3. Migrate clients gradually
4. Update documentation
5. Monitor performance improvements

## Future Enhancements

### Planned Features
- **Authentication Middleware**: Proper gRPC authentication interceptors
- **Streaming Support**: Real-time event updates
- **Pagination**: Efficient handling of large datasets
- **Rate Limiting**: Prevent abuse of services
- **Metrics**: Monitoring and observability
- **Load Balancing**: Distribute load across multiple instances

### Advanced Features
- **Bidirectional Streaming**: Real-time communication
- **Health Checks**: Service health monitoring
- **Interceptors**: Request/response middleware
- **Metadata**: Custom metadata handling
- **Deadlines**: Request timeout handling

---

This documentation provides a complete guide to implementing and using gRPC services in the Udemy Go Final Project. For additional questions or support, refer to the main README.md file or create an issue in the repository.

**Note:** This document serves as a reference guide. All implementation details, code examples, and working samples can be found in the referenced files throughout the project.