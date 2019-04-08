#!/bin/bash

dep ensure

find handlers -name main.go -type f \
 | xargs -n 1 dirname \
 | xargs -n 1 -I@ bash -c "cd ./@ && CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o main . && pwd"


