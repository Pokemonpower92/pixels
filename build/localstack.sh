#!/bin/bash

# Create the network
docker network create collage

# Create the db container
docker run -d \
  --name collage_db \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=collage \
  -v db:/var/lib/postgresql/data \
  --network collage \
  postgres:13

# Create the rabbitmq container
docker run -d \
  --name collage_rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  --network collage \
  rabbitmq:3-management

# Create the imagesetparser container
docker build -t imagesetparser:latest -f ./build/Dockerfile .
docker run \
    --name imagesetparser \
    --env-file ./.env.docker \
    --network collage \
    imagesetparser:latest

