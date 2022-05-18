#/bin/sh
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
  gateway
  transfer
  msg
  push
  sdk_server
  open_im_demo
)
#
version=v2.0.10
cd ../script/; ./build_all_service.sh
cd ../deploy_k8s/

for i in  ${service[*]}
do
  mv ../bin/open_im_${i} ./${i}/
done

echo "move success"

echo "start to build images"

for i in ${service[*]}
do
	echo "start to build images" $i
	cd $i
	image="openim/${i}:$version"
	docker build -t $image . -f ./${i}.Dockerfile
	echo "build ${dockerfile} success"
	docker push $image
	echo "push ${image} success "
	cd ..
done

