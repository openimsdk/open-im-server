set output_dir=%~dp0..\_output\bin\platforms\windows
go build -o %output_dir%\api.exe ../cmd/openim-api/main.go
go build -o %output_dir%\auth.exe ../cmd/openim-rpc/openim-rpc-auth/main.go
go build -o %output_dir%\conversation.exe ../cmd/openim-rpc/openim-rpc-conversation/main.go
go build -o %output_dir%\friend.exe ../cmd/openim-rpc/openim-rpc-friend/main.go
go build -o %output_dir%\group.exe ../cmd/openim-rpc/openim-rpc-group/main.go
go build -o %output_dir%\msg.exe ../cmd/openim-rpc/openim-rpc-msg/main.go
go build -o %output_dir%\third.exe ../cmd/openim-rpc/openim-rpc-third/main.go
go build -o %output_dir%\user.exe ../cmd/openim-rpc/openim-rpc-user/main.go
go build -o %output_dir%\push.exe ../cmd/openim-push/main.go
go build -o %output_dir%\msgtransfer.exe ../cmd/openim-msgtransfer/main.go
go build -o %output_dir%\msggateway.exe ../cmd/openim-msggateway/main.go