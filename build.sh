#!/bin/bash

# コンテナ名とイメージ名を指定
CONTAINER_NAME="go-auth"
IMAGE_NAME="go-auth-img"

# 既存のコンテナを停止して削除
echo "===== Stopping and removing existing container... ====="
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
echo ""

# 既存のイメージを削除
echo "===== Removing existing image... ====="
docker rmi $IMAGE_NAME
echo ""

# docker-compose.ymlから新たにイメージとコンテナを生成
echo "===== Building and starting new container from docker-compose.yml... ====="
docker-compose up --build -d
echo ""

echo "===== Done ======"