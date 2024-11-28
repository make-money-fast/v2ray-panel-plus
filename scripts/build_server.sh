#!/bin/bash

dir=$(dirname $1)

cp -r ./scripts/systemd $dir/
cp ./scripts/v2ray_server_install.sh $dir/