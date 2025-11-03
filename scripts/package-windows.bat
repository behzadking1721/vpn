@echo off
setlocal enabledelayedexpansion

echo Packaging VPN Client for Windows...

REM Get version from git or use default
for /f "tokens=*" %%a in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%a
if "%VERSION%"=="" set VERSION=dev

echo Packaging version: %VERSION%

REM Check if dist directory exists
if not exist "..\dist" (
    echo Error: Build directory 'dist' does not exist. Run build script first.
    exit /b 1
)

REM Change to project root directory
cd /d %~dp0..

REM Update version in Inno Setup script
echo Updating version in installer script...
echo #define MyAppName "VPN Client" > temp.iss
echo #define MyAppVersion "%VERSION%" >> temp.iss
echo #define MyAppPublisher "My Company" >> temp.iss
echo #define MyAppURL "https://www.mycompany.com/" >> temp.iss
echo #define MyAppExeName "vpn-client.exe" >> temp.iss
more +5 scripts\installer.iss >> temp.iss
move /y temp.iss scripts\installer.iss >nul

REM Check if Inno Setup compiler is available
set INNO_SETUP_FOUND=0
where ISCC >nul 2>&1
if %errorlevel% equ 0 (
    set INNO_SETUP_FOUND=1
) else (
    REM Check common installation paths
    if exist "C:\Program Files (x86)\Inno Setup 6\ISCC.exe" (
        set PATH=%PATH%;"C:\Program Files (x86)\Inno Setup 6\"
        set INNO_SETUP_FOUND=1
    )
    if exist "C:\Program Files\Inno Setup 6\ISCC.exe" (
        set PATH=%PATH%;"C:\Program Files\Inno Setup 6\"
        set INNO_SETUP_FOUND=1
    )
)

if %INNO_SETUP_FOUND% equ 0 (
    echo Warning: Inno Setup compiler not found. Please install Inno Setup to create installer.
    echo Creating ZIP archive instead...
    goto create_zip
)

REM Create installer using Inno Setup
echo Creating Windows installer with Inno Setup...
ISCC scripts\installer.iss

if exist "Output\vpn-client-setup.exe" (
    mkdir packages 2>nul
    move "Output\vpn-client-setup.exe" "packages\vpn-client-windows-%VERSION%-installer.exe"
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