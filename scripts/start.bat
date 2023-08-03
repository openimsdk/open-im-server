cd %~p0../_output/bin/platforms/windows
start api.exe -p 10002
start auth.exe -p 10060
start conversation.exe -p 10080
start friend.exe -p 10020
start group.exe -p 10050
start msg.exe -p 10030
start msggateway.exe -p 10040 -w 10001
start msgtransfer.exe
start third.exe -p 10090
start push.exe -p 10070
start user.exe -p 10010