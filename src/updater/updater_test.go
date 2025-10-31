package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdaterCreation(t *testing.T) {
	updater, err := NewUpdater("1.0.0", "http://example.com")
	if err != nil {
		t.Errorf("Failed to create updater: %v", err)
	}

	if updater == nil {
		t.Error("Expected updater to be created")
	}
}

func TestUpdaterWithInvalidVersion(t *testing.T) {
	_, err := NewUpdater("invalid-version", "http://example.com")
	if err == nil {
		t.Error("Expected error when creating updater with invalid version")
	}
}

func TestCheckForUpdateNoUpdate(t *testing.T) {
	// Create a mock server that returns the current version
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := ReleaseInfo{
			Version: "1.0.0",
			URL:     "http://example.com/vpn-client.exe",
		}
		json.NewEncoder(w).Encode(release)
	}))
	defer mockServer.Close()

	updater, err := NewUpdater("1.0.0", mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to create updater: %v", err)
	}

	release, err := updater.CheckForUpdate()
	if err != nil {
		t.Errorf("CheckForUpdate failed: %v", err)
	}

	if release != nil {
		t.Error("Expected no update to be available")
	}
}

func TestCheckForUpdateWithUpdate(t *testing.T) {
	// Create a mock server that returns a newer version
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := ReleaseInfo{
			Version: "1.0.1",
			URL:     "http://example.com/vpn-client.exe",
		}
		json.NewEncoder(w).Encode(release)
	}))
	defer mockServer.Close()

	updater, err := NewUpdater("1.0.0", mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to create updater: %v", err)
	}

	release, err := updater.CheckForUpdate()
	if err != nil {
		t.Errorf("CheckForUpdate failed: %v", err)
	}

	if release == nil {
		t.Error("Expected an update to be available")
	}

	if release.Version != "1.0.1" {
		t.Errorf("Expected version 1.0.1, got %s", release.Version)
	}
}
