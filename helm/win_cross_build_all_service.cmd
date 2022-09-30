SET ROOT=%cd%
mkdir %ROOT%\..\bin\
set GOARCH=amd64
go env -w GOARCH=amd64
set GOOS=linux
cd ..\cmd\open_im_api\ && go build -ldflags="-w -s" && move open_im_api %ROOT%\..\bin\
cd ..\..\cmd\open_im_cms_api\ && go build -ldflags="-w -s" && move open_im_cms_api %ROOT%\..\bin\
cd ..\..\cmd\open_im_demo\ && go build -ldflags="-w -s" && move open_im_demo %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_gateway\ && go build -ldflags="-w -s" && move open_im_msg_gateway %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_transfer\ && go build -ldflags="-w -s" && move open_im_msg_transfer %ROOT%\..\bin\
cd ..\..\cmd\open_im_push\ && go build -ldflags="-w -s" && move open_im_push %ROOT%\..\bin\
cd ..\..\cmd\rpc\open_im_admin_cms\&& go build -ldflags="-w -s" && move open_im_admin_cms %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_auth\&& go build -ldflags="-w -s" && move open_im_auth %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_cache\&& go build -ldflags="-w -s" && move open_im_cache %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_conversation\&& go build -ldflags="-w -s" && move open_im_conversation %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_friend\&& go build -ldflags="-w -s" && move open_im_friend %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_group\&& go build -ldflags="-w -s" && move open_im_group %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_message_cms\&& go build -ldflags="-w -s" && move open_im_message_cms %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_msg\&& go build -ldflags="-w -s" && move open_im_msg %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_office\&& go build -ldflags="-w -s" && move open_im_office %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_organization\&& go build -ldflags="-w -s" && move open_im_organization %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_statistics\&& go build -ldflags="-w -s" && move open_im_statistics %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_user\&& go build -ldflags="-w -s" && move open_im_user %ROOT%\..\bin\
cd ..\..\..\cmd\Open-IM-SDK-Core\ws_wrapper\cmd\&& go build -ldflags="-w -s" open_im_sdk_server.go && move open_im_sdk_server %ROOT%\..\bin\
go env -w GOARCH=amd64
go env -w GOOS=windows
cd %ROOT%