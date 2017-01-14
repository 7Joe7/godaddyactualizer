#!/usr/bin/env bash

set -x
set -e

IP_ADDRESS=192.168.1.214
LOGIN=root

env GOOS=linux go build

ssh -p 2202 $LOGIN@$IP_ADDRESS "systemctl stop godaddyactualizer.service; rm -rf /usr/local/src/godaddyactualizer; rm -rf /etc/systemd/system/godaddyactualizer.service"
scp -P 2202 ./godaddyactualizer $LOGIN@$IP_ADDRESS:/usr/local/src/godaddyactualizer
scp -P 2202 ./config.json $LOGIN@$IP_ADDRESS:/usr/local/src/config.json
scp -P 2202 ./godaddyactualizer.service $LOGIN@$IP_ADDRESS:/etc/systemd/system/godaddyactualizer.service
ssh -p 2202 $LOGIN@$IP_ADDRESS "systemctl daemon-reload; systemctl start godaddyactualizer.service"