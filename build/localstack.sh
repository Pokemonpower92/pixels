#!/bin/bash

# Create the db container
docker run -d \
  --name collage_db \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=posgres \
  -e POSTGRES_DB=collage \
  -v db:/var/lib/postgresql/data \
  --network collage \
  postgres:13

# Create the redis container
docker run -d \
  --name collage_redis \
  -p 6379:6379 \
  -v ./data:/data \
  --network collage \
  redis:6.2-alpine \
  redis-server --loglevel warning --appendonly yes --requirepass redispass

# Create the rabbitmq container
docker run -d \
  --name collage_rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  --network collage \
  rabbitmq:3-management

# Create the imageset service container
docker run --env-file ./.env --network collage imagesetservice:latest
