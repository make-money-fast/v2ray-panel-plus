#!/bin/bash

chmod +x v2ray
echo 'install v2ray to /usr/local/bin'
cp ./v2ray /usr/local/bin

echo 'install system service'
cp -r ./systemd /etc/

sudo systemctl daemon-reload
sudo systemctl start v2ray
sudo systemctl enable v2ray

echo 'commands:'
echo '启动 v2ray: systemctl start v2ray'
echo '停止 v2ray: systemctl stop v2ray'
echo '查看 v2ray: systemctl status v2ray'
echo ''

