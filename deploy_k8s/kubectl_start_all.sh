service=(
  #api service file
  api
  cms_api
  #rpc service file
  user
  friend
  group
  auth
  admin_cms
  message_cms
  statistics
  office
  organization
  conversation
  cache
  msg_gateway
  transfer
  msg
  push
  sdk_server
  demo
)

for i in ${service[*]}
do
  kubectl -n openim apply -f ./${i}/deployment.yaml
done

