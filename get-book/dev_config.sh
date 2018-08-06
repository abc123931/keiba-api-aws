#!/bin/sh

ifconfig en0 | sed -nEe $'s/^[ \t]+inet[ \t]+([0-9.]+).*$/\\1/p' > ip.txt
VAR=`cat < ip.txt`
echo "DYNAMO_ENPOINT="$VAR > .env
rm ip.txt
