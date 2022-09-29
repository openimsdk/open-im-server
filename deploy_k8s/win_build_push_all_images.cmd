SET ROOT=%cd%
@echo off
echo "ensure run Open-IM-Server/scipt/win_build_all_service.cmd first"
echo "ensure run Open-IM-Server/scipt/win_build_all_service.cmd first"
echo "ensure run Open-IM-Server/scipt/win_build_all_service.cmd first"
echo "ensure run Open-IM-Server/scipt/win_build_all_service.cmd first"
echo "ensure run Open-IM-Server/scipt/win_build_all_service.cmd first"
set version=v2.2.0
for %%I in (
  api
  cms_api
  user
  friend
  group
  auth
  admin_cms
  office
  organization
  conversation
  cache
  msg_gateway
  msg_transfer
  msg
  push
  sdk_server
  demo
) do  docker build -t openim/%%I:%version% --build-arg SER_NAME=%%I.exe ../ -f temp.Dockerfile
cd %ROOT%