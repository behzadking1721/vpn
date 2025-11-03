@echo off
setlocal enabledelayedexpansion

echo Packaging VPN Client for Windows...

REM Get version from git or use default
for /f "tokens=*" %%a in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%a
if "%VERSION%"=="" set VERSION=dev

echo Packaging version: %VERSION%

REM Check if dist directory exists
if not exist "dist" (
    echo Error: Build directory 'dist' does not exist. Run build script first.
    exit /b 1
)

REM Update version in Inno Setup script
echo Updating version in installer script...
powershell -Command "(Get-Content installer.iss) -replace '#define MyAppVersion .*', '#define MyAppVersion \"%VERSION%\"' | Set-Content installer.iss"

REM Check if Inno Setup compiler is available
where ISCC >nul 2>&1
if %errorlevel% neq 0 (
    echo Warning: Inno Setup compiler not found. Please install Inno Setup to create installer.
    echo Creating ZIP archive instead...
    goto create_zip
)

REM Create installer using Inno Setup
echo Creating Windows installer with Inno Setup...
ISCC installer.iss

if exist "vpn-client-setup.exe" (
    mkdir packages 2>nul
    move "vpn-client-setup.exe" "packages\vpn-client-windows-%VERSION%-installer.exe"
    echo Windows installer created successfully.
) else (
    echo Failed to create Windows installer.
)

:create_zip
REM Create ZIP archives for all Windows platforms
echo Creating ZIP archives...
mkdir packages 2>nul

if exist "dist\windows_amd64" (
    powershell -Command "Compress-Archive -Path 'dist\windows_amd64\*' -DestinationPath 'packages\vpn-client-windows-amd64-%VERSION%.zip' -Force"
    echo Created ZIP for windows_amd64
)

if exist "dist\windows_386" (
    powershell -Command "Compress-Archive -Path 'dist\windows_386\*' -DestinationPath 'packages\vpn-client-windows-386-%VERSION%.zip' -Force"
    echo Created ZIP for windows_386
)

if exist "dist\windows_arm64" (
    powershell -Command "Compress-Archive -Path 'dist\windows_arm64\*' -DestinationPath 'packages\vpn-client-windows-arm64-%VERSION%.zip' -Force"
    echo Created ZIP for windows_arm64
)

echo Windows packaging completed.
echo Packages are located in the packages\ directory.