@echo off
SETLOCAL EnableDelayedExpansion

SET "BIN_DIR=%~dp0..\_output\bin\platforms\windows\amd64"

SET "CONFIG_DIR=%~dp0..\config"

cd "%BIN_DIR%"

FOR %%f IN ("%BIN_DIR%\*.exe") DO (
    echo Starting %%~nf...
    start cmd /k "%%~f -i 0 -c "%CONFIG_DIR%" & echo Press any key to close this window... & pause>nul"
    echo %%~nf started.
)

echo All binaries in the directory have been started.
