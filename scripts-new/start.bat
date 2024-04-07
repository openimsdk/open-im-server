@echo off
SETLOCAL EnableDelayedExpansion

SET "BIN_DIR=%~dp0..\_output\bin\platforms\windows\amd64"

SET "CONFIG_DIR=%~dp0..\config"

cd "%BIN_DIR%"

FOR %%f IN ("%BIN_DIR%\*.exe") DO (
    echo Starting %%~nf...
    echo Command: start "" "%%~nf.exe" -i 0 -c "%CONFIG_DIR%"
    start "" "%%~nf.exe" -i 0 -c "%CONFIG_DIR%"
    echo %%~nf started.
)

echo All binaries in the directory have been started.
