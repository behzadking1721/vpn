@echo off
echo Running basic VPN Client tests...
echo.

echo Testing version command...
.\vpn-client.exe --version > version_output.txt
findstr /C:"VPN Client v1.0.0" version_output.txt >nul
if %errorlevel% equ 0 (
    echo PASS: Version command works correctly
) else (
    echo FAIL: Version command did not return expected output
    type version_output.txt
)

echo.
echo Testing help command...
.\vpn-client.exe --help > help_output.txt
findstr /C:"Usage: vpn-client [options]" help_output.txt >nul
if %errorlevel% equ 0 (
    echo PASS: Help command works correctly
) else (
    echo FAIL: Help command did not return expected output
    type help_output.txt
)

echo.
echo Cleaning up...
del version_output.txt help_output.txt

echo.
echo Basic tests completed.
pause