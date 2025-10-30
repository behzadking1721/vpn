package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"github.com/gorilla/mux"
)

func TestDashboardAPI(t *testing.T) {
	// Create mock managers
	configManager := managers.NewConfigManager("./test_config.json")
	serverManager := managers.NewServerManager()
	connManager := managers.NewConnectionManager()
	dataManager := serverManager.GetDataManager()
	
	// Create dashboard API
	dashboardAPI := NewDashboardAPI(serverManager, connManager, configManager, dataManager)
	
	// Create router
	router := mux.NewRouter()
	dashboardAPI.RegisterRoutes(router)
	
	// Test connection status endpoint
	t.Run("GetConnectionStatus", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/dashboard/status", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		
		// Check the response body is JSON
		var response ConnectionStatusResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
	})
	
	// Test data usage endpoint
	t.Run("GetDataUsage", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/dashboard/usage", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		
		// Check the response body is JSON
		var response DataUsageResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
	})
	
	// Test servers status endpoint
	t.Run("GetServersStatus", func(t *testing.T) {
		// Add a test server
		server := core.Server{
			ID:       utils.GenerateID(),
			Name:     "Test Server",
			Host:     "example.com",
			Port:     443,
			Protocol: core.ProtocolVMess,
			Enabled:  true,
		}
		
		if err := serverManager.AddServer(server); err != nil {
			t.Fatalf("failed to add server: %v", err)
		}
		
		req, err := http.NewRequest("GET", "/api/dashboard/servers", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		
		// Check the response body is JSON
		var response ServerStatusResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
		
		if len(response.Servers) != 1 {
			t.Errorf("expected 1 server, got %d", len(response.Servers))
		}
	})
	
	// Test statistics endpoint
	t.Run("GetStatistics", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/dashboard/stats", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		
		// Check the response body is JSON
		var response StatisticsResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
	})
	
	// Test recent logs endpoint
	t.Run("GetRecentLogs", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/dashboard/logs", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		
		// Check the response body is JSON
		var response RecentLogsResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
	})
}