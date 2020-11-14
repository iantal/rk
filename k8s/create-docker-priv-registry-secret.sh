#!/bin/bash

DOCKER_REGISTRY_SERVER=docker.io
DOCKER_USER=micky0000
DOCKER_EMAIL=antal.micky@gmail.com

echo -n 'Password:'
 
read -s DOCKER_PASSWORD

kubectl create secret docker-registry dockerregistrykey \
  --docker-server=$DOCKER_REGISTRY_SERVER \
  --docker-username=$DOCKER_USER \
  --docker-password=$DOCKER_PASSWORD \
  --docker-email=$DOCKER_EMAIL