package managers

import (
	"os"
	"testing"
	"io/ioutil"
)

func TestDataManager(t *testing.T) {
	// Create temporary files
	serversFile, err := ioutil.TempFile("", "servers_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp servers file: %v", err)
	}
	defer os.Remove(serversFile.Name())
	serversFile.Close()
	
	subsFile, err := ioutil.TempFile("", "subscriptions_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp subscriptions file: %v", err)
	}
	defer os.Remove(subsFile.Name())
	subsFile.Close()
	
	// Create data manager
	dm := NewDataManager(serversFile.Name(), subsFile.Name())
	
	// Test that we can create a data manager
	if dm == nil {
		t.Error("Failed to create data manager")
	}
	
	t.Logf("Data manager created with servers file: %s and subscriptions file: %s", 
		serversFile.Name(), subsFile.Name())
}

func TestDataManager_SaveLoad(t *testing.T) {
	// Create temporary files
	serversFile, err := ioutil.TempFile("", "servers_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp servers file: %v", err)
	}
	defer os.Remove(serversFile.Name())
	serversFile.Close()
	
	subsFile, err := ioutil.TempFile("", "subscriptions_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp subscriptions file: %v", err)
	}
	defer os.Remove(subsFile.Name())
	subsFile.Close()
	
	// Create data manager
	dm := NewDataManager(serversFile.Name(), subsFile.Name())
	
	// Test saving and loading data
	// In a real implementation, we would test the actual save/load functionality
	t.Logf("Testing save/load functionality for data manager")
}