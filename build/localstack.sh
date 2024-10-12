#!/usr/bin/bash

build_flag='false'
start_flag='false'
clean_flag='false'

network="collage"

# Images 
rmq_image="rabbitmq:3-management"
postgres_image="postgres:13"
imagesetparser_image="imagesetparser:latest"

images="$rmq_image $postgres_image $imagesetparser_image"

# Containers
rmq='collage_rabbitmq'
db='collage_db'
imagesetparser='imagesetparser'

containers="$rmq $db $imagesetparser"

print_usage() {
  printf "Usage: ..."
}

while getopts 'bsc' flag; do
  case "${flag}" in
    b) build_flag='true' ;;
    s) start_flag='true' ;;
    c) clean_flag='true' ;;
    *) print_usage
       exit 1 ;;
  esac
done

clean ()
{
    docker stop $containers
    docker rm $containers
    docker image rm $images
    docker network rm $network
}

if [[ $build_flag == 'true' ]]; then
    clean

    # Create the network
    docker network create $network

    # Create the db container
    docker run -d \
      --name $db \
      -p 5432:5432 \
      -e POSTGRES_USER=postgres \
      -e POSTGRES_PASSWORD=postgres \
      -e POSTGRES_DB=collage \
      -v db:/var/lib/postgresql/data \
      --network $network \
      $postgres_image

    # Create the rabbitmq container
    docker run -d \
      --name $rmq \
      -p 5672:5672 \
      -p 15672:15672 \
      -e RABBITMQ_DEFAULT_USER=guest \
      -e RABBITMQ_DEFAULT_PASS=guest \
      --network $network \
      $rmq_image
    
     # Create the imagesetparser container
    docker build -t $imagesetparser_image -f ./build/Dockerfile .
    docker run \
        --name $imagesetparser \
        --env-file ./.env.docker \
        --network $network \
        $imagesetparser_image
fi

if [[ $start_flag == 'true' ]]; then
    docker start $rmq
    docker start $db
    sleep 5
    docker start $imagesetparser
fi

if [[ $clean_flag == 'true' ]]; then
   clean 
fi

