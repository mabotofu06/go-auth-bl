#!/bin/sh

image=go-auth:postgre-db

docker rm -f $image

docker build -t $image .

docker volume rm -f go-auth-db_db-store

docker rm -f go-auth-db
docker compose up -d