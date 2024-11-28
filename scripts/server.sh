#!/bin/bash

# 设置 GitHub 仓库的拥有者和名称
OWNER="make-money-fast"
REPO="v2ray-panel-plus"

# 获取最新 Release 的版本号
tag=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | grep '"tag_name":' | cut -d '"' -f 4)

# 检查是否成功获取到版本号
if [ -z "$tag" ]; then
    echo "未获取到版本号"
    exit 1
fi

echo "最新版本: $tag"

filename="v2ray-panel-plus_$tag_linux_amd64.tar.gz"

# 下载指定文件（根据需要修改文件名和下载链接）
# 这里假设下载的是一个名为 "example-linux-amd64" 的文件
download_url="https://github.com/$OWNER/$REPO/releases/download/$tag/$filename"
curl -L -o "$filename" "$download_url"

tar -zxvf $filename
