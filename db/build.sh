#!/bin/sh

image=go-auth-db:postgre-db
container=go-auth-db

# 古いコンテナの削除
echo "===== Stopping and removing existing container... ====="
docker stop
docker rm -f $container

# 古いイメージの削除
docker rmi -f $image

# Dockerイメージのビルド
docker build -t $image .

# 古いボリュームの削除
docker volume rm -f db_db-store

# Docker Composeを使用してコンテナを起動
docker-compose up -d