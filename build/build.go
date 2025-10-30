package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// BuildConfig holds configuration for building the application
type BuildConfig struct {
	AppName     string
	Version     string
	Platforms   []Platform
	BuildDir    string
	DistDir     string
	IncludeDirs []string
	IncludeFiles []string
}

// Platform represents a target platform for building
type Platform struct {
	OS        string
	Arch      string
	Extension string
}

func main() {
	// Define build configuration
	config := BuildConfig{
		AppName: "vpn-client",
		Version: "1.0.0",
		Platforms: []Platform{
			{OS: "windows", Arch: "amd64", Extension: ".exe"},
			{OS: "windows", Arch: "386", Extension: ".exe"},
			{OS: "linux", Arch: "amd64", Extension: ""},
			{OS: "linux", Arch: "arm64", Extension: ""},
			{OS: "linux", Arch: "arm", Extension: ""},
			{OS: "darwin", Arch: "amd64", Extension: ""},
			{OS: "darwin", Arch: "arm64", Extension: ""},
		},
		BuildDir: "build",
		DistDir:  "dist",
		IncludeDirs: []string{
			"ui/desktop",
			"config",
			"resources",
		},
		IncludeFiles: []string{
			"LICENSE",
		},
	}

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "help" {
		printHelp()
		return
	}

	// Build for all platforms or specific ones
	if len(args) > 0 {
		// Build for specific platform
		targetOS := args[0]
		targetArch := "amd64"
		if len(args) > 1 {
			targetArch = args[1]
		}
		buildForPlatform(config, targetOS, targetArch)
	} else {
		// Build for all platforms
		buildAllPlatforms(config)
	}

	fmt.Println("Build process completed!")
}

func printHelp() {
	helpText := `
VPN Client Build Tool

Usage:
  go run build.go [os] [arch]    Build for specific platform
  go run build.go               Build for all platforms
  go run build.go help          Show this help

Examples:
  go run build.go windows amd64  Build for Windows 64-bit
  go run build.go linux arm64    Build for Linux ARM64
  go run build.go                Build for all platforms

Supported platforms:
  windows/amd64, windows/386
  linux/amd64, linux/arm64, linux/arm
  darwin/amd64, darwin/arm64
`
	fmt.Println(helpText)
}

func buildAllPlatforms(config BuildConfig) {
	fmt.Println("Building for all platforms...")

	for _, platform := range config.Platforms {
		fmt.Printf("Building for %s/%s...\n", platform.OS, platform.Arch)
		err := buildForPlatform(config, platform.OS, platform.Arch)
		if err != nil {
			log.Printf("Failed to build for %s/%s: %v", platform.OS, platform.Arch, err)
		}
	}
}

func buildForPlatform(config BuildConfig, targetOS, targetArch string) error {
	// Set environment variables for cross-compilation
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOOS=%s", targetOS))
	env = append(env, fmt.Sprintf("GOARCH=%s", targetArch))

	// Create platform-specific directory
	platformDir := filepath.Join(config.DistDir, fmt.Sprintf("%s-%s", targetOS, targetArch))
	err := os.MkdirAll(platformDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create platform directory: %v", err)
	}

	// Determine executable name
	execName := config.AppName
	if targetOS == "windows" {
		execName += ".exe"
	}

	// Build the application
	fmt.Printf("Building %s for %s/%s...\n", execName, targetOS, targetArch)
	cmd := exec.Command("go", "build", "-o", filepath.Join(platformDir, execName), "src/main.go")
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to build application: %v", err)
	}

	// Copy additional files and directories
	err = copyAdditionalFiles(config, platformDir)
	if err != nil {
		return fmt.Errorf("failed to copy additional files: %v", err)
	}

	// Create package based on platform
	switch targetOS {
	case "windows":
		err = createWindowsPackage(config, platformDir, targetArch)
	case "linux":
		err = createLinuxPackage(config, platformDir, targetArch)
	case "darwin":
		err = createMacOSPackage(config, platformDir, targetArch)
	}

	if err != nil {
		return fmt.Errorf("failed to create package: %v", err)
	}

	fmt.Printf("Successfully built for %s/%s\n", targetOS, targetArch)
	return nil
}

func copyAdditionalFiles(config BuildConfig, platformDir string) error {
	// Copy directories
	for _, dir := range config.IncludeDirs {
		destDir := filepath.Join(platformDir, filepath.Base(dir))
		err := copyDir(dir, destDir)
		if err != nil {
			return fmt.Errorf("failed to copy directory %s: %v", dir, err)
		}
	}

	// Copy files
	for _, file := range config.IncludeFiles {
		destFile := filepath.Join(platformDir, filepath.Base(file))
		err := copyFile(file, destFile)
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %v", file, err)
		}
	}

	// Create necessary directories
	dirs := []string{"data", "logs", "config"}
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(platformDir, dir), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

func createWindowsPackage(config BuildConfig, platformDir, targetArch string) error {
	// For Windows, we'll create a simple zip archive
	// In a real implementation, you might want to create an MSI or NSIS installer
	zipName := fmt.Sprintf("%s_%s_windows_%s.zip", config.AppName, config.Version, targetArch)
	zipPath := filepath.Join(config.DistDir, zipName)

	fmt.Printf("Creating Windows package: %s\n", zipName)

	// Create zip command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Use PowerShell on Windows
		cmd = exec.Command("powershell", "-Command", 
			fmt.Sprintf("Compress-Archive -Path '%s\\*' -DestinationPath '%s' -Force", 
				strings.ReplaceAll(platformDir, "/", "\\"), 
				strings.ReplaceAll(zipPath, "/", "\\")))
	} else {
		// Use zip command on Unix-like systems
		cmd = exec.Command("zip", "-r", zipPath, ".")
		cmd.Dir = platformDir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createLinuxPackage(config BuildConfig, platformDir, targetArch string) error {
	// For Linux, we'll create a tar.gz archive
	// In a real implementation, you might want to create DEB or RPM packages
	tarName := fmt.Sprintf("%s_%s_linux_%s.tar.gz", config.AppName, config.Version, targetArch)
	tarPath := filepath.Join(config.DistDir, tarName)

	fmt.Printf("Creating Linux package: %s\n", tarName)

	// Create tar.gz command
	cmd := exec.Command("tar", "-czf", tarPath, "-C", filepath.Dir(platformDir), filepath.Base(platformDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createMacOSPackage(config BuildConfig, platformDir, targetArch string) error {
	// For macOS, we'll create a tar.gz archive
	// In a real implementation, you might want to create a DMG or PKG installer
	tarName := fmt.Sprintf("%s_%s_macos_%s.tar.gz", config.AppName, config.Version, targetArch)
	tarPath := filepath.Join(config.DistDir, tarName)

	fmt.Printf("Creating macOS package: %s\n", tarName)

	// Create tar.gz command
	cmd := exec.Command("tar", "-czf", tarPath, "-C", filepath.Dir(platformDir), filepath.Base(platformDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file contents
	_, err = dstFile.ReadFrom(srcFile)
	return err
}