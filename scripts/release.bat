@echo off
setlocal enabledelayedexpansion

REM Release script for VPN Client (Windows)
REM Usage: scripts\release.bat [version]

REM Colors for output (using PowerShell for colored output)
set "GREEN=Write-Host"
set "YELLOW=Write-Host"
set "RED=Write-Host"

REM Function to print status
set "PRINT_STATUS=%GREEN% '[STATUS]' -NoNewline; Write-Host "

REM Function to print warning
set "PRINT_WARNING=%YELLOW% '[WARNING]' -NoNewline; Write-Host "

REM Function to print error
set "PRINT_ERROR=%RED% '[ERROR]' -NoNewline; Write-Host "

REM Check if version is provided
if "%1"=="" (
    powershell -Command "%PRINT_ERROR% 'Version not specified'"
    echo Usage: %0 ^<version^>
    echo Example: %0 1.2.0
    exit /b 1
)

set "VERSION=%1"

REM Validate version format (semantic versioning)
echo %VERSION% | findstr /R "^[0-9]*\.[0-9]*\.[0-9]*$" >nul
if errorlevel 1 (
    powershell -Command "%PRINT_ERROR% 'Invalid version format. Use semantic versioning (e.g., 1.2.0)'"
    exit /b 1
)

powershell -Command "%PRINT_STATUS% 'Starting release process for version %VERSION%'"

REM Check git status
git status --porcelain | findstr /R "." >nul
if not errorlevel 1 (
    powershell -Command "%PRINT_ERROR% 'Git working directory is not clean. Please commit or stash changes.'"
    exit /b 1
)

REM Get current version
set "CURRENT_VERSION=0.0.0"
if exist versioninfo.json (
    for /f "tokens=4 delims=:\" %%a in ('findstr "FileVersion" versioninfo.json') do (
        set "CURRENT_VERSION=%%a"
        REM Clean up the version string
        set "CURRENT_VERSION=!CURRENT_VERSION:"=!"
        set "CURRENT_VERSION=!CURRENT_VERSION: =!"
        set "CURRENT_VERSION=!CURRENT_VERSION:,=!"
    )
)

powershell -Command "%PRINT_STATUS% 'Current version: !CURRENT_VERSION!'"

REM Update version in files
powershell -Command "%PRINT_STATUS% 'Updating version from !CURRENT_VERSION! to %VERSION%'"

REM Update versioninfo.json
if exist versioninfo.json (
    powershell -Command "(Get-Content versioninfo.json) -replace 'FileVersion.*:!CURRENT_VERSION!', 'FileVersion\": \"%VERSION%\"' | Set-Content versioninfo.json"
    powershell -Command "(Get-Content versioninfo.json) -replace 'ProductVersion.*:!CURRENT_VERSION!', 'ProductVersion\": \"%VERSION%\"' | Set-Content versioninfo.json"
)

REM Update installer.nsi
if exist installer.nsi (
    powershell -Command "(Get-Content installer.nsi) -replace '!define APP_VERSION.*!CURRENT_VERSION!', '!define APP_VERSION \"%VERSION%\"' | Set-Content installer.nsi"
)

REM Update or create changelog
if exist CHANGELOG.md (
    REM Add new version to changelog (this is a simplified approach)
    powershell -Command "Set-Content -Path 'temp_changelog.md' -Value '## [%VERSION%] - %date%'" 
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value ''"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '### Added'"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '- '"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value ''"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '### Changed'"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '- '"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value ''"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '### Fixed'"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value '- '"
    powershell -Command "Add-Content -Path 'temp_changelog.md' -Value ''"
    powershell -Command "Get-Content CHANGELOG.md | Select-Object -Skip 1 | Add-Content -Path 'temp_changelog.md'"
    move /y temp_changelog.md CHANGELOG.md >nul
) else (
    REM Create changelog
    echo # Changelog > CHANGELOG.md
    echo. >> CHANGELOG.md
    echo All notable changes to this project will be documented in this file. >> CHANGELOG.md
    echo. >> CHANGELOG.md
    echo The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), >> CHANGELOG.md
    echo and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). >> CHANGELOG.md
    echo. >> CHANGELOG.md
    echo ## [%VERSION%] - %date% >> CHANGELOG.md
    echo. >> CHANGELOG.md
    echo ### Added >> CHANGELOG.md
    echo - >> CHANGELOG.md
    echo. >> CHANGELOG.md
    echo ### Changed >> CHANGELOG.md
    echo - >> CHANGELOG.md
    echo. >> CHANGELOG.md
    echo ### Fixed >> CHANGELOG.md
    echo - >> CHANGELOG.md
    echo. >> CHANGELOG.md
)

REM Build application
powershell -Command "%PRINT_STATUS% 'Building application version %VERSION%'"

REM Clean previous builds
if exist build rmdir /s /q build
if exist dist rmdir /s /q dist
if exist bin rmdir /s /q bin

REM Build using make (if make is available)
make build-all 2>nul
if errorlevel 1 (
    REM Fallback to manual build if make is not available
    powershell -Command "%PRINT_WARNING% 'Make not found, using manual build process'"
    
    REM Create bin directory
    mkdir bin 2>nul
    
    REM Build for Windows
    go build -o bin\vpn-client-windows-amd64.exe src\main.go
    
    REM Build for Linux
    set GOOS=linux
    set GOARCH=amd64
    go build -o bin\vpn-client-linux-amd64 src\main.go
    
    REM Build for macOS
    set GOOS=darwin
    set GOARCH=amd64
    go build -o bin\vpn-client-darwin-amd64 src\main.go
    
    REM Reset GOOS and GOARCH
    set GOOS=
    set GOARCH=
)

REM Create distribution packages
powershell -Command "%PRINT_STATUS% 'Creating distribution packages'"

REM Create releases directory
mkdir releases 2>nul
mkdir releases\v%VERSION% 2>nul

REM Copy built binaries to releases directory
if exist bin (
    xcopy /E /I /Y bin releases\v%VERSION%\ >nul
)

REM Create git tag
powershell -Command "%PRINT_STATUS% 'Creating git tag v%VERSION%'"
git add .
git commit -m "Release version %VERSION%"
git tag -a "v%VERSION%" -m "Release version %VERSION%"

powershell -Command "%PRINT_STATUS% 'Release v%VERSION% created successfully!'"
powershell -Command "%PRINT_STATUS% 'Release files are located in releases/v%VERSION%/'"
powershell -Command "%PRINT_STATUS% 'Don't forget to push the changes and tag:'"
echo   git push origin main
echo   git push origin v%VERSION%

endlocal