#!/usr/bin/env bash

r=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 20 | head -1)
echo "swagger-${r}"
