#!/bin/sh

ifconfig en0 | sed -nEe $'s/^[ \t]+inet[ \t]+([0-9.]+).*$/\\1/p' > ip.txt
VAR=`cat < ip.txt`
echo "DYNAMO_ENPOINT=http://"$VAR":8000" > .env_dev
rm ip.txt
cp .env_dev ./gethorsename/.env
cp .env_dev ./getracename/.env
rm .env_dev
