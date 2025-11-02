package managers

import (
	"testing"
)

func TestNewConnectionManager(t *testing.T) {
	cm := NewConnectionManager()
	if cm == nil {
		t.Error("Failed to create ConnectionManager")
		return
	}

	status := cm.GetStatus()
	if status != Disconnected {
		t.Errorf("Expected initial status to be Disconnected, got %v", status)
	}
}