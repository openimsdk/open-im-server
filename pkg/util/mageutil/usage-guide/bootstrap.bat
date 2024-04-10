@echo off
SETLOCAL

mage -version >nul 2>&1
IF %ERRORLEVEL% EQU 0 (
    echo Mage is already installed.
    GOTO DOWNLOAD
)

go version >nul 2>&1
IF NOT %ERRORLEVEL% EQU 0 (
    echo Go is not installed. Please install Go and try again.
    exit /b 1
)

echo Installing Mage...
go install github.com/magefile/mage@latest

mage -version >nul 2>&1
IF NOT %ERRORLEVEL% EQU 0 (
    echo Mage installation failed.
    echo Please ensure that %GOPATH%/bin is in your PATH.
    exit /b 1
)

echo Mage installed successfully.

:DOWNLOAD
go mod download

ENDLOCAL
