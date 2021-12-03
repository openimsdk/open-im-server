#!/usr/bin/env bash
source ./style_info.cfg

echo "docker-compose ps..................................."
docker-compose ps


echo  -e "check OpenIM result............................."
i=1
t=5
while [ $i -le $t ]
do
        p=`awk 'BEGIN{printf "%.2f%\n",('$i'/'$t')*100}'`
        echo -e ${GREEN_PREFIX} "=> $p"${COLOR_SUFFIX}
        sleep 5
        let i++
done

./check_all.sh
