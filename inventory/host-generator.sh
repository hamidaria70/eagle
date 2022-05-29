#!/bin/bash

set -e

SERVERS=$(grep hostname ~/.ssh/config | cut -d " " -f 4 )

echo [servers] > ./inventory/hosts

for IP in ${SERVERS[@]};do	
	echo $IP >> ./inventory/hosts
done

