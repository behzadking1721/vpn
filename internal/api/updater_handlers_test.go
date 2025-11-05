package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/internal/notifications"
	"vpnclient/internal/stats"
	"vpnclient/internal/updater"
	"vpnclient/mocks"
)

// MockUpdater is a mock implementation of the Updater
type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) GetStatus() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func (m *MockUpdater) SetEnabled(enabled bool) {
	m.Called(enabled)
}

func (m *MockUpdater) SetInterval(interval time.Duration) {
	m.Called(interval)
}

func (m *MockUpdater) UpdateSubscriptions() error {
	args := m.Called()
	return args.Error(0)
}

// TestGetUpdaterStatus tests the getUpdaterStatus handler
func TestGetUpdaterStatus(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Mock expectations
	expectedStatus := map[string]interface{}{
		"enabled":  true,
		"interval": "24h0m0s",
	}
	mockUpdater.On("GetStatus").Return(expectedStatus)

	// Create request
	req, err := http.NewRequest("GET", "/api/updater/status", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, true, response["enabled"])
	assert.Equal(t, "24h0m0s", response["interval"])
	
	mockUpdater.AssertExpectations(t)
}

// TestSetUpdaterConfigEnabled tests setting updater enabled status
func TestSetUpdaterConfigEnabled(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Mock expectations
	mockUpdater.On("SetEnabled", true).Return()

	// Create request body
	config := map[string]interface{}{
		"enabled": true,
	}
	
	body, err := json.Marshal(config)
	assert.NoError(t, err)

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/config", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Updater configuration updated successfully", response["message"])
	
	mockUpdater.AssertExpectations(t)
}

// TestSetUpdaterConfigInterval tests setting updater interval
func TestSetUpdaterConfigInterval(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Mock expectations
	mockUpdater.On("SetInterval", 6*time.Hour).Return()

	// Create request body
	config := map[string]interface{}{
		"interval": "6h",
	}
	
	body, err := json.Marshal(config)
	assert.NoError(t, err)

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/config", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Updater configuration updated successfully", response["message"])
	
	mockUpdater.AssertExpectations(t)
}

// TestSetUpdaterConfigInvalidInterval tests setting invalid updater interval
func TestSetUpdaterConfigInvalidInterval(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Create request body with invalid interval
	config := map[string]interface{}{
		"interval": "invalid",
	}
	
	body, err := json.Marshal(config)
	assert.NoError(t, err)

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/config", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "Invalid interval format: time: invalid duration \"invalid\"", response["error"])
}

// TestSetUpdaterConfigInvalidJSON tests setting updater config with invalid JSON
func TestSetUpdaterConfigInvalidJSON(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Create request with invalid JSON
	body := []byte("{invalid json}")

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/config", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Contains(t, response["error"], "Invalid JSON")
}

// TestTriggerUpdate tests triggering an update
func TestTriggerUpdate(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Mock expectations
	mockUpdater.On("UpdateSubscriptions").Return(nil)

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/update", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Subscription update triggered successfully", response["message"])
	
	mockUpdater.AssertExpectations(t)
}

// TestTriggerUpdateError tests triggering an update that fails
func TestTriggerUpdateError(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Mock expectations
	mockUpdater.On("UpdateSubscriptions").Return(assert.AnError)

	// Create request
	req, err := http.NewRequest("POST", "/api/updater/update", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "Failed to update subscriptions: assert.AnError general error for testing", response["error"])
	
	mockUpdater.AssertExpectations(t)
}

// TestTriggerUpdateWrongMethod tests triggering an update with wrong HTTP method
func TestTriggerUpdateWrongMethod(t *testing.T) {
	// Setup
	logger := &logging.Logger{}
	mockUpdater := &MockUpdater{}
	
	server := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		updater: mockUpdater,
	}
	
	server.setupRoutes()

	// Create request with wrong method
	req, err := http.NewRequest("GET", "/api/updater/update", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Test
	server.router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "Method not allowed", response["error"])
}