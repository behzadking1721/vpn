package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	
	"github.com/blang/semver/v4"
	update "github.com/inconshreveable/go-update"
)

// ReleaseInfo holds information about a release
type ReleaseInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	Checksum    string `json:"checksum"`
	ReleaseDate string `json:"release_date"`
	Notes       string `json:"notes"`
}

// Updater handles application updates
type Updater struct {
	currentVersion semver.Version
	apiURL         string
	logger         *utils.Logger
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string, apiURL string) (*Updater, error) {
	version, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %v", err)
	}

	logger := utils.NewLogger(&utils.LoggerConfig{
		Level: utils.LogLevelInfo,
		File:  "./logs/updater.log",
	})

	return &Updater{
		currentVersion: version,
		apiURL:         apiURL,
		logger:         logger,
	}, nil
}

// CheckForUpdate checks if there is a newer version available
func (u *Updater) CheckForUpdate() (*ReleaseInfo, error) {
	u.logger.Info("Checking for updates...")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(u.apiURL)
	if err != nil {
		u.logger.Error("Failed to check for updates: %v", err)
		return nil, fmt.Errorf("failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		u.logger.Error("Failed to check for updates: HTTP %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to check for updates: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		u.logger.Error("Failed to parse release info: %v", err)
		return nil, fmt.Errorf("failed to parse release info: %v", err)
	}

	// Parse the release version
	releaseVersion, err := semver.ParseTolerant(release.Version)
	if err != nil {
		u.logger.Error("Invalid release version format: %v", err)
		return nil, fmt.Errorf("invalid release version format: %v", err)
	}

	// Compare versions
	if releaseVersion.GT(u.currentVersion) {
		u.logger.Info("New version available: %s", release.Version)
		return &release, nil
	}

	u.logger.Info("Application is up to date")
	return nil, nil
}

// DownloadUpdate downloads the update file
func (u *Updater) DownloadUpdate(release *ReleaseInfo) (string, error) {
	u.logger.Info("Downloading update from %s", release.URL)

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Get(release.URL)
	if err != nil {
		u.logger.Error("Failed to download update: %v", err)
		return "", fmt.Errorf("failed to download update: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		u.logger.Error("Failed to download update: HTTP %d", resp.StatusCode)
		return "", fmt.Errorf("failed to download update: HTTP %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "vpn-client-update-*.exe")
	if err != nil {
		u.logger.Error("Failed to create temporary file: %v", err)
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	u.logger.Info("Saving update to %s", tmpFile.Name())

	// Copy the downloaded content to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		u.logger.Error("Failed to save update: %v", err)
		return "", fmt.Errorf("failed to save update: %v", err)
	}

	return tmpFile.Name(), nil
}

// ApplyUpdate applies the downloaded update
func (u *Updater) ApplyUpdate(updateFilePath string) error {
	u.logger.Info("Applying update from %s", updateFilePath)

	// Open the update file
	updateFile, err := os.Open(updateFilePath)
	if err != nil {
		u.logger.Error("Failed to open update file: %v", err)
		return fmt.Errorf("failed to open update file: %v", err)
	}
	defer updateFile.Close()

	// Apply the update
	err = update.Apply(updateFile, update.Options{})
	if err != nil {
		u.logger.Error("Failed to apply update: %v", err)
		return fmt.Errorf("failed to apply update: %v", err)
	}

	u.logger.Info("Update applied successfully")
	return nil
}

// Cleanup removes temporary files
func (u *Updater) Cleanup(tempFiles ...string) {
	for _, file := range tempFiles {
		if file != "" {
			if err := os.Remove(file); err != nil {
				u.logger.Warn("Failed to remove temporary file %s: %v", file, err)
			} else {
				u.logger.Info("Removed temporary file %s", file)
			}
		}
	}
}

// GetCurrentVersion returns the current application version
func (u *Updater) GetCurrentVersion() semver.Version {
	return u.currentVersion
}

// GetExecutablePath returns the path of the current executable
func GetExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return filepath.Clean(execPath), nil
}

// RestartApplication restarts the application
func RestartApplication() error {
	execPath, err := GetExecutablePath()
	if err != nil {
		return err
	}

	// In a real implementation, you would start the new process and exit the current one
	// This is a simplified version for demonstration purposes
	fmt.Printf("Application should restart from: %s\n", execPath)
	return nil
}