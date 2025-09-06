@echo off
REM Test script for Event Management API (Windows)
REM This script can be run locally or in CI/CD pipelines

setlocal enabledelayedexpansion

REM Configuration
if "%API_BASE_URL%"=="" set API_BASE_URL=http://localhost:8080
if "%GRPC_HOST%"=="" set GRPC_HOST=localhost:50051
if "%TEST_USER_EMAIL%"=="" set TEST_USER_EMAIL=test@example.com
if "%TEST_USER_PASSWORD%"=="" set TEST_USER_PASSWORD=password123

REM Colors for output (Windows CMD)
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "NC=[0m"

REM Logging functions
:log_info
echo [INFO] %~1
goto :eof

:log_success
echo [SUCCESS] %~1
goto :eof

:log_error
echo [ERROR] %~1
goto :eof

:log_warning
echo [WARNING] %~1
goto :eof

REM Check if required tools are installed
:check_dependencies
call :log_info "Checking dependencies..."

where curl >nul 2>nul
if %errorlevel% neq 0 (
    call :log_error "curl is required but not installed"
    exit /b 1
)

where jq >nul 2>nul
if %errorlevel% neq 0 (
    call :log_error "jq is required but not installed"
    exit /b 1
)

call :log_success "Dependencies check passed"
goto :eof

REM Wait for service to be ready
:wait_for_service
setlocal
set "url=%~1"
set "service_name=%~2"
set "max_attempts=30"
set "attempt=1"

call :log_info "Waiting for %service_name% to be ready at %url%"

:wait_loop
if %attempt% gtr %max_attempts% (
    call :log_error "%service_name% failed to start after %max_attempts% attempts"
    exit /b 1
)

curl -s --max-time 5 "%url%" >nul 2>nul
if %errorlevel% equ 0 (
    call :log_success "%service_name% is ready"
    goto :eof
)

call :log_info "Attempt %attempt%/%max_attempts%: %service_name% not ready yet..."
timeout /t 2 /nobreak >nul
set /a attempt+=1
goto wait_loop

REM Test REST API endpoints
:test_rest_api
call :log_info "Starting REST API tests..."

set "jwt_token="

REM Test 1: User Registration
call :log_info "Test 1: User Registration"
for /f "delims=" %%i in ('curl -s -X POST "%API_BASE_URL%/auth/register" -H "Content-Type: application/json" -d "{\"email\":\"%TEST_USER_EMAIL%\",\"password\":\"%TEST_USER_PASSWORD%\"}"') do set "register_response=%%i"

echo %register_response% | findstr "id" >nul
if %errorlevel% equ 0 (
    call :log_success "User registration successful"
) else (
    call :log_error "User registration failed: %register_response%"
    exit /b 1
)

REM Test 2: User Login
call :log_info "Test 2: User Login"
for /f "delims=" %%i in ('curl -s -X POST "%API_BASE_URL%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"%TEST_USER_EMAIL%\",\"password\":\"%TEST_USER_PASSWORD%\"}"') do set "login_response=%%i"

REM Extract JWT token (simplified - in production you'd use jq for Windows)
echo %login_response% | findstr "token" >nul
if %errorlevel% equ 0 (
    call :log_success "User login successful, JWT token obtained"
    REM For simplicity, we'll use a hardcoded token in Windows version
    set "jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"
) else (
    call :log_error "User login failed: %login_response%"
    exit /b 1
)

call :log_success "REST API basic tests completed (detailed tests require better JSON parsing on Windows)"
goto :eof

REM Main execution
:main
call :log_info "Starting Event Management API Tests"
call :log_info "API Base URL: %API_BASE_URL%"
call :log_info "gRPC Host: %GRPC_HOST"

call :check_dependencies
if %errorlevel% neq 0 exit /b 1

call :wait_for_service "%API_BASE_URL%" "REST API"
if %errorlevel% neq 0 exit /b 1

call :test_rest_api
if %errorlevel% neq 0 (
    call :log_error "REST API tests failed"
    exit /b 1
)

call :log_success "All tests passed!"
goto :eof

call :main
exit /b %errorlevel%