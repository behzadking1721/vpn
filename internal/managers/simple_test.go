package managers

import (
	"testing"
)

// TestSimple checks basic functionality
func TestSimple(t *testing.T) {
	// Simple test to verify the testing framework works
	result := 1 + 1
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}
