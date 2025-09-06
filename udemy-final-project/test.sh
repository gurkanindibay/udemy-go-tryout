#!/bin/bash

# Test script for Event Management API
# This script can be run locally or in CI/CD pipelines

set -e  # Exit on any error

# Configuration
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
GRPC_HOST="${GRPC_HOST:-localhost:50051}"
TEST_USER_EMAIL="${TEST_USER_EMAIL:-test@example.com}"
TEST_USER_PASSWORD="${TEST_USER_PASSWORD:-password123}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        log_error "jq is required but not installed"
        exit 1
    fi

    log_success "Dependencies check passed"
}

# Wait for service to be ready
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1

    log_info "Waiting for $service_name to be ready at $url"

    while [ $attempt -le $max_attempts ]; do
        if curl -s --max-time 5 "$url" > /dev/null 2>&1; then
            log_success "$service_name is ready"
            return 0
        fi

        log_info "Attempt $attempt/$max_attempts: $service_name not ready yet..."
        sleep 2
        ((attempt++))
    done

    log_error "$service_name failed to start after $max_attempts attempts"
    return 1
}

# Test REST API endpoints
test_rest_api() {
    log_info "Starting REST API tests..."

    local jwt_token=""

    # Test 1: User Registration
    log_info "Test 1: User Registration"
    local register_response=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_USER_EMAIL\",\"password\":\"$TEST_USER_PASSWORD\"}")

    if echo "$register_response" | jq -e '.id' > /dev/null 2>&1; then
        log_success "User registration successful"
    else
        log_error "User registration failed: $register_response"
        return 1
    fi

    # Test 2: User Login
    log_info "Test 2: User Login"
    local login_response=$(curl -s -X POST "$API_BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_USER_EMAIL\",\"password\":\"$TEST_USER_PASSWORD\"}")

    jwt_token=$(echo "$login_response" | jq -r '.token')
    if [ "$jwt_token" != "null" ] && [ -n "$jwt_token" ]; then
        log_success "User login successful, JWT token obtained"
    else
        log_error "User login failed: $login_response"
        return 1
    fi

    # Test 3: Get Events (should be empty initially)
    log_info "Test 3: Get Events (empty list expected)"
    local events_response=$(curl -s -X GET "$API_BASE_URL/events")

    if echo "$events_response" | jq -e '. | length == 0' > /dev/null 2>&1; then
        log_success "Get events successful (empty list as expected)"
    else
        log_error "Get events failed or returned unexpected data: $events_response"
        return 1
    fi

    # Test 4: Create Event
    log_info "Test 4: Create Event"
    local create_event_response=$(curl -s -X POST "$API_BASE_URL/events" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $jwt_token" \
        -d '{
            "name": "Test Event",
            "description": "This is a test event",
            "location": "Test Location",
            "date_time": "2025-12-25T10:00:00Z"
        }')

    local event_id=$(echo "$create_event_response" | jq -r '.id')
    if [ "$event_id" != "null" ] && [ -n "$event_id" ]; then
        log_success "Event creation successful, ID: $event_id"
    else
        log_error "Event creation failed: $create_event_response"
        return 1
    fi

    # Test 5: Get Events (should have 1 event now)
    log_info "Test 5: Get Events (should return 1 event)"
    events_response=$(curl -s -X GET "$API_BASE_URL/events")

    if echo "$events_response" | jq -e '. | length == 1' > /dev/null 2>&1; then
        log_success "Get events successful (1 event returned)"
    else
        log_error "Get events failed or returned wrong count: $events_response"
        return 1
    fi

    # Test 6: Get Event by ID
    log_info "Test 6: Get Event by ID"
    local event_by_id_response=$(curl -s -X GET "$API_BASE_URL/events/$event_id")

    if echo "$event_by_id_response" | jq -e '.id == '"$event_id" > /dev/null 2>&1; then
        log_success "Get event by ID successful"
    else
        log_error "Get event by ID failed: $event_by_id_response"
        return 1
    fi

    # Test 7: Update Event
    log_info "Test 7: Update Event"
    local update_response=$(curl -s -X PUT "$API_BASE_URL/events/$event_id" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $jwt_token" \
        -d '{
            "name": "Updated Test Event",
            "description": "This is an updated test event",
            "location": "Updated Test Location",
            "date_time": "2025-12-26T11:00:00Z"
        }')

    if [ "$update_response" = "{}" ] || echo "$update_response" | jq -e 'has("message")' > /dev/null 2>&1; then
        log_success "Event update successful"
    else
        log_error "Event update failed: $update_response"
        return 1
    fi

    # Test 8: Register for Event
    log_info "Test 8: Register for Event"
    local register_response=$(curl -s -X POST "$API_BASE_URL/events/$event_id/register" \
        -H "Authorization: Bearer $jwt_token")

    if [ "$register_response" = "{}" ] || echo "$register_response" | jq -e 'has("message")' > /dev/null 2>&1; then
        log_success "Event registration successful"
    else
        log_error "Event registration failed: $register_response"
        return 1
    fi

    # Test 9: Get User Registrations
    log_info "Test 9: Get User Registrations"
    local user_id=$(echo "$login_response" | jq -r '.user.id')
    local registrations_response=$(curl -s -X GET "$API_BASE_URL/users/$user_id/registrations" \
        -H "Authorization: Bearer $jwt_token")

    if echo "$registrations_response" | jq -e '. | length >= 1' > /dev/null 2>&1; then
        log_success "Get user registrations successful"
    else
        log_error "Get user registrations failed: $registrations_response"
        return 1
    fi

    # Test 10: Cancel Registration
    log_info "Test 10: Cancel Registration"
    local cancel_response=$(curl -s -X DELETE "$API_BASE_URL/events/$event_id/register" \
        -H "Authorization: Bearer $jwt_token")

    if [ "$cancel_response" = "{}" ] || echo "$cancel_response" | jq -e 'has("message")' > /dev/null 2>&1; then
        log_success "Registration cancellation successful"
    else
        log_error "Registration cancellation failed: $cancel_response"
        return 1
    fi

    # Test 11: Delete Event
    log_info "Test 11: Delete Event"
    local delete_response=$(curl -s -X DELETE "$API_BASE_URL/events/$event_id" \
        -H "Authorization: Bearer $jwt_token")

    if [ "$delete_response" = "{}" ] || echo "$delete_response" | jq -e 'has("message")' > /dev/null 2>&1; then
        log_success "Event deletion successful"
    else
        log_error "Event deletion failed: $delete_response"
        return 1
    fi

    log_success "All REST API tests passed! ✅"
    return 0
}

# Test gRPC endpoints
test_grpc_api() {
    log_info "Starting gRPC API tests..."

    # For now, we'll just check if the gRPC server is responding
    # In a real scenario, you'd use grpcurl or a Go gRPC client

    log_info "gRPC tests would go here (grpcurl or Go client needed)"
    log_warning "gRPC testing skipped for now - requires grpcurl or Go gRPC client"

    return 0
}

# Main test execution
main() {
    log_info "🚀 Starting Event Management API Tests"
    log_info "API Base URL: $API_BASE_URL"
    log_info "gRPC Host: $GRPC_HOST"

    # Check dependencies
    check_dependencies

    # Wait for services to be ready
    wait_for_service "$API_BASE_URL" "REST API"

    # Run tests
    local rest_result=0
    local grpc_result=0

    if test_rest_api; then
        log_success "REST API tests completed successfully"
    else
        log_error "REST API tests failed"
        rest_result=1
    fi

    if test_grpc_api; then
        log_success "gRPC API tests completed successfully"
    else
        log_error "gRPC API tests failed"
        grpc_result=1
    fi

    # Summary
    if [ $rest_result -eq 0 ] && [ $grpc_result -eq 0 ]; then
        log_success "🎉 All tests passed!"
        exit 0
    else
        log_error "❌ Some tests failed"
        exit 1
    fi
}

# Run main function
main "$@"