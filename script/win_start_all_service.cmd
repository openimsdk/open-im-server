SET ROOT=%cd%
cd %ROOT%\..\bin\
start cmd /C .\open_im_api.exe -port 10002
start cmd /C .\open_im_user.exe -port 10110
start cmd /C .\open_im_friend.exe -port 10120
start cmd /C .\open_im_group.exe -port 10150
start cmd /C .\open_im_auth.exe -port 10160
start cmd /C .\open_im_message_cms.exe -port 10190
start cmd /C .\open_im_msg.exe -port 10130
start cmd /C .\open_im_conversation.exe -port 10230
start cmd /C .\open_im_push.exe -port 10170
start cmd /C .\open_im_msg_transfer.exe
start cmd /C .\sdkws_server.exe -openIM_api_port 10002 -openIM_ws_port 10001 -sdkws_port 10003 -openIM_log_level 6
start cmd /C .\open_im_msg_gateway.exe -rpc_port 10140 -ws_port 10001
cd %ROOT%