#!/bin/bash

# Docker compose stub for the 3 local minio clusters.
# Add the necessary docker-compose commands (e.g. 'up -d', 'down', ...).

cd $(dirname $0)
cd ../deployment

docker compose -f minio.docker-compose.yaml $@
