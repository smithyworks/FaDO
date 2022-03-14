#!/bin/bash

# Docker compose stub for the 3 mock FaaS servers.
# Add the necessary docker-compose commands (e.g. 'up -d', 'down', ...).

cd $(dirname $0)
cd ../deployment

docker compose -f faas.docker-compose.yaml $@
