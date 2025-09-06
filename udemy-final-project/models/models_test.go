package models

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gurkanindibay/udemy-rest-api/db"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Save(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user
	user := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Test Save method
	err := user.Save()
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestGetUserByEmail(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user
	testUser := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	err := testUser.Save()
	require.NoError(t, err)

	// Test GetUserByEmail
	user, err := GetUserByEmail("test@example.com")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestVerifyUserCredentials(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user
	testUser := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	err := testUser.Save()
	require.NoError(t, err)

	// Test valid credentials
	user, err := VerifyUserCredentials("test@example.com", "testpassword")
	require.NoError(t, err)
	assert.NotNil(t, user)

	// Test invalid credentials
	user, err = VerifyUserCredentials("test@example.com", "wrongpassword")
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestEvent_Save(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user first
	testUser := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	err := testUser.Save()
	require.NoError(t, err)

	// Create test event
	event := Event{
		Name:        "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour), // Future date
		UserId:      testUser.ID,
	}

	// Test Save method
	err = event.Save()
	require.NoError(t, err)
	assert.NotZero(t, event.ID)
}

func TestGetAllEvents(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user
	testUser := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	err := testUser.Save()
	require.NoError(t, err)

	// Create test event
	event := Event{
		Name:        "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      testUser.ID,
	}
	err = event.Save()
	require.NoError(t, err)

	// Test GetAllEvents
	events, err := GetAllEvents()
	require.NoError(t, err)
	assert.NotEmpty(t, events)

	found := false
	for _, e := range events {
		if e.ID == event.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created event not found in GetAllEvents result")
}

func TestGetEventByID(t *testing.T) {
	// Setup test database
	setupTestDB(t)

	// Create test user
	testUser := User{
		Email:    "test@example.com",
		Password: "testpassword",
	}
	err := testUser.Save()
	require.NoError(t, err)

	// Create test event
	event := Event{
		Name:        "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserId:      testUser.ID,
	}
	err = event.Save()
	require.NoError(t, err)

	// Test GetEventByID
	retrievedEvent, err := GetEventByID(fmt.Sprintf("%d", event.ID))
	require.NoError(t, err)
	assert.NotNil(t, retrievedEvent)
	assert.Equal(t, event.ID, retrievedEvent.ID)
}

// Helper function to setup test database
func setupTestDB(t *testing.T) *sql.DB {
	// Initialize database connection if not already done
	if db.GetDB() == nil {
		db.InitDB()
	}

	// Clean up test data before each test
	cleanupTestData(t)

	return db.GetDB()
}

// Helper function to clean up test data
func cleanupTestData(t *testing.T) {
	testDB := db.GetDB()

	// Delete test data in correct order to respect foreign key constraints
	// 1. Delete registrations first (no dependencies)
	_, err := testDB.Exec("DELETE FROM registrations WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%')")
	if err != nil {
		t.Logf("Warning: Failed to clean up registrations: %v", err)
	}

	// 2. Delete events (depends on users)
	_, err = testDB.Exec("DELETE FROM events WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'test%')")
	if err != nil {
		t.Logf("Warning: Failed to clean up events: %v", err)
	}

	// 3. Delete users last (after events are deleted)
	_, err = testDB.Exec("DELETE FROM users WHERE email LIKE 'test%'")
	if err != nil {
		t.Logf("Warning: Failed to clean up users: %v", err)
	}
}
