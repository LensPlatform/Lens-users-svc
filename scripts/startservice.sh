#!/usr/bin/env bash
echo ">>>>>>>>>>>>>>>LENS<<<<<<<<<<<<<<<"
echo "stopping any running containers"
docker-compose stop
echo "removing stopped container"
docker-compose rm -f
echo "pulling latest containers from docker hub"
docker-compose pull
echo "spinning up services"
docker-compose up