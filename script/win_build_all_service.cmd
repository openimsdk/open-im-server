SET ROOT=%cd%
mkdir %ROOT%\..\bin\
cd ..\cmd\open_im_api\ && go build -ldflags="-w -s" && move open_im_api.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_cms_api\ && go build -ldflags="-w -s" && move open_im_cms_api.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_demo\ && go build -ldflags="-w -s" && move open_im_demo.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_gateway\ && go build -ldflags="-w -s" && move open_im_msg_gateway.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_transfer\ && go build -ldflags="-w -s" && move open_im_msg_transfer.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_push\ && go build -ldflags="-w -s" && move open_im_push.exe %ROOT%\..\bin\
cd ..\..\cmd\rpc\open_im_admin_cms\&& go build -ldflags="-w -s" && move open_im_admin_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_auth\&& go build -ldflags="-w -s" && move open_im_auth.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_cache\&& go build -ldflags="-w -s" && move open_im_cache.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_conversation\&& go build -ldflags="-w -s" && move open_im_conversation.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_friend\&& go build -ldflags="-w -s" && move open_im_friend.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_group\&& go build -ldflags="-w -s" && move open_im_group.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_message_cms\&& go build -ldflags="-w -s" && move open_im_message_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_msg\&& go build -ldflags="-w -s" && move open_im_msg.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_office\&& go build -ldflags="-w -s" && move open_im_office.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_organization\&& go build -ldflags="-w -s" && move open_im_organization.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_statistics\&& go build -ldflags="-w -s" && move open_im_statistics.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_user\&& go build -ldflags="-w -s" && move open_im_user.exe %ROOT%\..\bin\
cd ..\..\..\cmd\Open-IM-SDK-Core\ws_wrapper\cmd\&& go build -ldflags="-w -s" open_im_sdk_server.go && move open_im_sdk_server.exe %ROOT%\..\bin\
cd %ROOT%