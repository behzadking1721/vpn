@echo off
setlocal

echo Building VPN Client for Windows...

REM Set variables
set APP_NAME=vpn-client
set BUILD_DIR=build
set DIST_DIR=dist

REM Create build directories
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
if not exist %DIST_DIR% mkdir %DIST_DIR%
if not exist logs mkdir logs
if not exist config mkdir config
if not exist data mkdir data

REM Generate version info
echo Generating version info...
goversioninfo -o %BUILD_DIR%\resource.syso versioninfo.json

REM Build the application
echo Building application...
go build -o %BUILD_DIR%\%APP_NAME%.exe -v src/main.go

REM Copy necessary files
echo Copying files...
copy /Y %BUILD_DIR%\%APP_NAME%.exe %DIST_DIR%\
xcopy /E /I /Y ui\desktop %DIST_DIR%\ui\
if not exist %DIST_DIR%\config mkdir %DIST_DIR%\config
if not exist %DIST_DIR%\logs mkdir %DIST_DIR%\logs
if not exist %DIST_DIR%\data mkdir %DIST_DIR%\data

REM Create installer with NSIS (if available)
if exist installer.nsi (
    echo Creating installer...
    makensis installer.nsi
    if exist %APP_NAME%-installer.exe (
        move /Y %APP_NAME%-installer.exe %DIST_DIR%\
        echo Installer created: %DIST_DIR%\%APP_NAME%-installer.exe
    )
)

echo.
echo Build completed successfully!
echo Executable: %DIST_DIR%\%APP_NAME%.exe
if exist %DIST_DIR%\%APP_NAME%-installer.exe (
    echo Installer: %DIST_DIR%\%APP_NAME%-installer.exe
)
echo.

pause