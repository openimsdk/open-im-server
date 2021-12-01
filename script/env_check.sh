#!/usr/bin/env bash
source ./style_info.cfg
echo -e "check environment......................................."

SYSTEM=`uname -s`
if [ $SYSTEM != "Linux" ] ; then
        echo -e ${RED_PREFIX}"Warning: Currently only Linux is supported"${COLOR_SUFFIX}
else
        echo -e ${GREEN_PREFIX} "Linux system is ok"${COLOR_SUFFIX}
fi

echo -e "check memory............................................"
available=`free -m | grep Mem | awk '{print $NF}'`
if [ $available -lt 2000 ] ; then
        echo -e ${RED_PREFIX}"Warning: Your memory not enough, available is: " "$available"m${COLOR_SUFFIX}"\c"
        echo -e ${RED_PREFIX}", must be greater than 2000m"${COLOR_SUFFIX}
else
        echo -e ${GREEN_PREFIX} "Memory is ok, available is: "$available"m${COLOR_SUFFIX}"
fi

