#!/usr/bin/env bash

# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${SCRIPTS_ROOT}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/path_info.sh
source $SCRIPTS_ROOT/lib/init.sh

cd $SCRIPTS_ROOT

echo -e "check time synchronize.................................."
t=`curl http://time.akamai.com/?iso -s`
t1=`date -d $t +%s`
t2=`date +%s`
let between=t2-t1
if [[ $between -gt 10 ]] || [[ $between -lt -10 ]]; then
  echo -e ${RED_PREFIX}"Warning: The difference between the iso time and the server's time is too large: "$between"s" ${COLOR_SUFFIX}
else
   echo -e ${GREEN_PREFIX} "ok: Server time is synchronized " ${COLOR_SUFFIX}
fi


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
