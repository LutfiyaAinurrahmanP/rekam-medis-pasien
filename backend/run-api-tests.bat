@echo off
REM Test runner script with Jest-like output for Windows
REM This script runs all API tests and provides a nice summary

echo.
echo ╔════════════════════════════════════════════════════════════════╗
echo ║                  API Test Suite Runner                        ║
echo ║                  (Jest-like Output)                            ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.

cd /d "%~dp0"
go test -v ./internal/routes/... -count=1

echo.
pause
