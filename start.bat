cd /d %~dp0
start api.exe -p 10002
start user.exe -p 10010
start friend.exe -p 10020
start group.exe -p 10050
start msg.exe -p 10030
start third.exe -p 10090
start conversation.exe -p 10080
start push.exe -p 10070
start auth.exe -p 10060
start msgtransfer.exe
start msggateway.exe -p 10040 -w 10001