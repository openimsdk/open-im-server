#!/usr/bin/env bash
source ./style_info.cfg

echo -e "check login user........................................"
user=`whoami`
if [ $user == "root" ] ; then
        echo -e ${GREEN_PREFIX} "ok: login user is root" ${COLOR_SUFFIX}
else
        echo -e ${RED_PREFIX}"Warning: The current user is not root "${COLOR_SUFFIX}
fi



echo -e "check docker............................................"
docker_running=`systemctl status docker | grep running |  grep active | wc -l`

docker_version=`docker-compose -v; docker -v`

if [ $docker_running -gt 0 ]; then
	echo -e ${GREEN_PREFIX} "ok: docker is running"   ${COLOR_SUFFIX}
	echo -e ${GREEN_PREFIX}  $docker_version ${COLOR_SUFFIX}

else
	echo -e ${RED_PREFIX}"docker not running"${COLOR_SUFFIX}
fi


echo -e "check environment......................................."
SYSTEM=`uname -s`
if [ $SYSTEM != "Linux" ] ; then
        echo -e ${RED_PREFIX}"Warning: Currently only Linux is supported"${COLOR_SUFFIX}
else
        echo -e ${GREEN_PREFIX} "ok: system is linux"${COLOR_SUFFIX}
fi

echo -e "check memory............................................"
available=`free -m | grep Mem | awk '{print $NF}'`
if [ $available -lt 2000 ] ; then
        echo -e ${RED_PREFIX}"Warning: Your memory not enough, available is: " "$available"m${COLOR_SUFFIX}"\c"
        echo -e ${RED_PREFIX}", must be greater than 2000m"${COLOR_SUFFIX}
else
        echo -e ${GREEN_PREFIX} "ok: available memory is: "$available"m${COLOR_SUFFIX}"
fi
