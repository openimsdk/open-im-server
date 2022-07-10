#/bin/sh
source ./path_info.cfg

# images version
version=v2.2.0
git pull
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

