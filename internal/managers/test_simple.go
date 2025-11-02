package managers

import (
	"testing"
)

// TestSimplest checks basic functionality
func TestSimplest(t *testing.T) {
	// Simple test to verify the testing framework works
	result := 2 + 2
	if result != 4 {
		t.Errorf("Expected 4, got %d", result)
	}
}