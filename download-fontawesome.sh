#!/bin/bash

# FontAwesome下载脚本
# 将FontAwesome 6.4.0资源下载到static目录

set -e

echo "Downloading FontAwesome 6.4.0..."

# 创建目录（使用webfonts以匹配CSS中的路径）
mkdir -p static/css static/webfonts

# 下载CSS文件
echo "Downloading CSS..."
curl -s -L -o static/css/font-awesome.min.css \
  https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css

# 下载字体文件
echo "Downloading fonts..."
curl -s -L -o static/webfonts/fa-solid-900.woff2 \
  https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/webfonts/fa-solid-900.woff2

curl -s -L -o static/webfonts/fa-regular-400.woff2 \
  https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/webfonts/fa-regular-400.woff2

curl -s -L -o static/webfonts/fa-brands-400.woff2 \
  https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/webfonts/fa-brands-400.woff2

echo "FontAwesome download completed!"
echo "Files are in static/ directory"