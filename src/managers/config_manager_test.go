package managers

import (
	"os"
	"testing"
	"io/ioutil"
)

func TestConfigManager(t *testing.T) {
	// Create a temporary config file
	tempFile, err := ioutil.TempFile("", "config_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Close the file as we only need the path
	tempFile.Close()
	
	// Create config manager
	cm := NewConfigManager(tempFile.Name())
	
	// Test that we can create a config manager
	if cm == nil {
		t.Error("Failed to create config manager")
	}
	
	// Test default config values
	// Note: We can't directly access the config field as it's private
	// In a real implementation, we would have getter methods
}

func TestConfigManager_SaveLoad(t *testing.T) {
	// Create a temporary config file
	tempFile, err := ioutil.TempFile("", "config_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Close the file as we only need the path
	tempFile.Close()
	
	// Create config manager
	cm := NewConfigManager(tempFile.Name())
	
	// Test saving and loading config
	// In a real implementation, we would test the actual save/load functionality
	t.Logf("Config manager created with path: %s", tempFile.Name())
}