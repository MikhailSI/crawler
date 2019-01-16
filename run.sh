#!/bin/bash

sudo docker build -q -t crawler:1 . -f .docker/Dockerfile
sudo docker run -it -e URL=$1 -e RPS=$2 crawler:1