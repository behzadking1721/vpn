package monitoring

import (
	"testing"
	"time"
)

func TestHealthManager(t *testing.T) {
	// Create a new health manager
	healthManager := NewHealthManager()

	// Test adding checkers
	t.Run("AddChecker", func(t *testing.T) {
		// Create a mock checker
		checker := &mockChecker{
			name: "test-checker",
			result: CheckResult{
				Status:    "ok",
				Message:   "Test checker is working",
				Timestamp: time.Now(),
			},
		}

		// Add the checker
		healthManager.AddChecker(checker)

		// Verify the checker was added
		addedChecker := healthManager.GetChecker("test-checker")
		if addedChecker == nil {
			t.Error("Failed to add checker")
		}
	})

	// Test removing checkers
	t.Run("RemoveChecker", func(t *testing.T) {
		// Add another mock checker
		checker := &mockChecker{
			name: "remove-test-checker",
			result: CheckResult{
				Status:    "ok",
				Message:   "Test checker for removal",
				Timestamp: time.Now(),
			},
		}

		healthManager.AddChecker(checker)

		// Verify the checker was added
		addedChecker := healthManager.GetChecker("remove-test-checker")
		if addedChecker == nil {
			t.Error("Failed to add checker for removal")
		}

		// Remove the checker
		healthManager.RemoveChecker("remove-test-checker")

		// Verify the checker was removed
		removedChecker := healthManager.GetChecker("remove-test-checker")
		if removedChecker != nil {
			t.Error("Failed to remove checker")
		}
	})

	// Test performing checks
	t.Run("PerformChecks", func(t *testing.T) {
		// Add a healthy checker
		healthyChecker := &mockChecker{
			name: "healthy-checker",
			result: CheckResult{
				Status:    "ok",
				Message:   "Healthy checker",
				Timestamp: time.Now(),
			},
		}

		// Add an unhealthy checker
		unhealthyChecker := &mockChecker{
			name: "unhealthy-checker",
			result: CheckResult{
				Status:    "error",
				Message:   "Unhealthy checker",
				Timestamp: time.Now(),
			},
		}

		healthManager.AddChecker(healthyChecker)
		healthManager.AddChecker(unhealthyChecker)

		// Perform checks
		result := healthManager.Check()

		// Verify the result
		if result.Status != "unhealthy" {
			t.Errorf("Expected overall status 'unhealthy', got '%s'", result.Status)
		}

		if len(result.Checks) != 3 { // 2 mock checkers + 1 system checker
			t.Errorf("Expected 3 checks, got %d", len(result.Checks))
		}

		// Verify individual check results
		healthyResult, ok := result.Checks["healthy-checker"]
		if !ok {
			t.Error("Healthy checker result not found")
		} else if healthyResult.Status != "ok" {
			t.Errorf("Expected healthy checker status 'ok', got '%s'", healthyResult.Status)
		}

		unhealthyResult, ok := result.Checks["unhealthy-checker"]
		if !ok {
			t.Error("Unhealthy checker result not found")
		} else if unhealthyResult.Status != "error" {
			t.Errorf("Expected unhealthy checker status 'error', got '%s'", unhealthyResult.Status)
		}
	})
}

func TestSystemHealthChecker(t *testing.T) {
	// Create a system health checker
	checker := NewSystemHealthChecker()

	// Perform the check
	result := checker.Check()

	// Verify the result
	if result.Status != "ok" {
		t.Errorf("Expected system health status 'ok', got '%s'", result.Status)
	}

	if result.Message == "" {
		t.Error("System health message is empty")
	}

	if result.Timestamp.IsZero() {
		t.Error("System health timestamp is zero")
	}
}

// mockChecker is a mock implementation of the HealthChecker interface for testing
type mockChecker struct {
	name   string
	result CheckResult
}

func (m *mockChecker) Check() CheckResult {
	return m.result
}

func (m *mockChecker) Name() string {
	return m.name
}
