#!/bin/bash

# Run the Caddy load balancer.

cd $(dirname $0)
cd ../deployment

docker run -it \
  -p 2019:2019 \
  -p 6000:6000 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  -v ${PWD}/load-balancer.Caddyfile:/etc/caddy/Caddyfile \
  --rm \
  caddy:latest
