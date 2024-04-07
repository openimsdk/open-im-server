@echo off
SETLOCAL EnableDelayedExpansion

SET "ROOT_DIR=%~dp0..\"
SET "OUTPUT_DIR=%ROOT_DIR%_output\bin\platforms\windows\amd64\"

IF NOT EXIST "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

call :findMainGo "%ROOT_DIR%cmd"

echo Compilation complete.
goto :eof

:findMainGo
FOR /R %1 %%d IN (.) DO (
    IF EXIST "%%d\main.go" (
        SET "DIR_PATH=%%d"
        SET "DIR_NAME=%%~nxd"

        echo Found main.go in %%d
        echo Compiling %%d...


        pushd "%%d"
        SET "GOOS=windows"
        SET "GOARCH=amd64"
        go build -o "!OUTPUT_DIR!!DIR_NAME!.exe" main.go
        if ERRORLEVEL 1 (
            echo Failed to compile %%d
            goto :eof
        )
        popd
    )
)
goto :eof
