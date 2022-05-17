#/bin/sh
#service_source_root=(
#  #api service file
#  ../cmd/open_im_api/
#  ../cmd/open_im_cms_api/
#  #rpc service file
#  ../cmd/rpc/open_im_user/
#  ../cmd/rpc/open_im_friend/
#  ../cmd/rpc/open_im_group/
#  ../cmd/rpc/open_im_auth/
#  ../cmd/rpc/open_im_admin_cms/
#  ../cmd/rpc/open_im_message_cms/
#  ../cmd/rpc/open_im_statistics/
#  ../cmd/rpc/open_im_office/
#  ../cmd/rpc/open_im_organization/
#  ../cmd/rpc/open_im_conversation/
#  ../cmd/rpc/open_im_cache/
#  ../cmd/open_im_msg_gateway/
#  ../cmd/open_im_msg_transfer/
#  ../cmd/rpc/open_im_msg/
#  ../cmd/open_im_push/
#  ../cmd/Open-IM-SDK-Core/
#  ../cmd/open_im_demo/
#)
#
dockerfile_list=pwd

docker build -t  $image . -f deploy.Dockerfile
