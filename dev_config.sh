#!/bin/sh

ifconfig en0 | sed -nEe $'s/^[ \t]+inet[ \t]+([0-9.]+).*$/\\1/p' > ip.txt
VAR=`cat < ip.txt`
echo "DYNAMO_ENPOINT=http://"$VAR":8000" > .env_dev
rm ip.txt
cp .env_dev ./gethorsename/.env
cp .env_dev ./getracename/.env
cp .env_dev ./getcourseresult/.env
echo "CHANNEL_SECRET=9c4b6ca4f431d5ce01167ae61eacab0c" >> .env_dev
echo "CHANNEL_TOKEN=AWczPVqdSmVXqVDuYq2eqo15dSnTtV3NFn6z8sYMDGgTWRx755i+5mbKdWk6m9uRW0DQaBSjHBHm7P/UFVrHVqKPe3x1od/F6FRFidRTVbnpc60n0992p" >> .env_dev
cp .env_dev ./line_bot_test/.env
rm .env_dev
