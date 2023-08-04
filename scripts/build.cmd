@echo off
set output_dir=%~dp0..\_output\bin\platforms\windows

set "rpc_apps=auth conversation friend group msg third user"
set "other_apps=api push msgtransfer msggateway"

for %%a in (%rpc_apps%) do (
    go build -o %output_dir%\%%a.exe ../cmd/openim-rpc/openim-rpc-%%a/main.go
)

for %%a in (%other_apps%) do (
    go build -o %output_dir%\%%a.exe ../cmd/openim-%%a/main.go
)
