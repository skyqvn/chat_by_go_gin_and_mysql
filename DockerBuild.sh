#!/bin/bash

docker build -t chat:v1 .
docker run --name chat -p 80:80 -d chat:v1
