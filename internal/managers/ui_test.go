package managers

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getUIPath determines the correct path to UI files based on the OS and file structure
func getUIPath() string {
	_, filename, _, _ := runtime.Caller(0)
	// Navigate from internal/managers to ui/desktop
	basePath := filepath.Join(filepath.Dir(filename), "..", "..", "ui", "desktop")
	return basePath
}

// TestUIElements tests that the UI HTML files contain expected elements
func TestUIElements(t *testing.T) {
	// Get the correct UI path
	uiPath := getUIPath()
	
	// Test index.html
	testUIFile(t, filepath.Join(uiPath, "index.html"), []string{
		"VPN Client",
		"connect-button",
		"serverSelect",
		"dataUsageChart",
		"fastestServerBtn",
		"testPingBtn",
		"bestServerBtn",
	})

	// Test dashboard.html
	testUIFile(t, filepath.Join(uiPath, "dashboard.html"), []string{
		"VPN Client - Dashboard",
		"dataUsageChart",
		"topServersList",
		"connectionHistory",
	})

	// Test alerts.html
	testUIFile(t, filepath.Join(uiPath, "alerts.html"), []string{
		"VPN Client - Alerts",
		"alertsList",
		"serverList",
	})
}

// testUIFile is a helper function to test a UI HTML file
func testUIFile(t *testing.T, filePath string, expectedElements []string) {
	// Read the HTML file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", filePath, err)
	}

	// Convert content to string
	htmlContent := string(content)

	// Check for expected elements
	for _, element := range expectedElements {
		if !strings.Contains(htmlContent, element) {
			t.Errorf("Expected element '%s' not found in %s", element, filePath)
		}
	}
}

// TestUIFunctionality tests UI functionality by simulating user interactions
func TestUIFunctionality(t *testing.T) {
	// This would typically involve a headless browser or similar tool
	// For now, we'll just test that the required JavaScript functions exist
	t.Skip("UI functionality tests require a headless browser environment")
}

// TestUINavigation tests that navigation between UI pages works correctly
func TestUINavigation(t *testing.T) {
	uiPath := getUIPath()
	
	pages := []string{
		filepath.Join(uiPath, "index.html"),
		filepath.Join(uiPath, "dashboard.html"),
		filepath.Join(uiPath, "alerts.html"),
	}

	// Check that each page has navigation links to other pages
	for _, page := range pages {
		testNavigationInPage(t, page, []string{
			"index.html",
			"dashboard.html",
			"alerts.html",
		})
	}
}

// testNavigationInPage is a helper function to test navigation in a specific page
func testNavigationInPage(t *testing.T, filePath string, expectedLinks []string) {
	// Read the HTML file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", filePath, err)
	}

	// Convert content to string
	htmlContent := string(content)

	// Check for expected navigation links
	for _, link := range expectedLinks {
		if !strings.Contains(htmlContent, link) {
			t.Errorf("Expected navigation link '%s' not found in %s", link, filePath)
		}
	}
}

// TestUITheme tests that UI theme functionality is present
func TestUITheme(t *testing.T) {
	uiPath := getUIPath()
	
	// Check that theme-related elements exist in the UI files
	testUIFile(t, filepath.Join(uiPath, "index.html"), []string{
		"themeToggle",
		"theme-light",
		"theme-dark",
	})

	testUIFile(t, filepath.Join(uiPath, "dashboard.html"), []string{
		"themeToggle",
		"theme-light",
		"theme-dark",
	})

	testUIFile(t, filepath.Join(uiPath, "alerts.html"), []string{
		"themeToggle",
		"theme-light",
		"theme-dark",
	})
}