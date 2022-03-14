#!/bin/bash

# Run the current server.

cd $(dirname $0)
cd ../server

go run main.go -c "../deployment/standalone.server-config.json" \
  --server-url http://localhost:9090 \
  --database postgres://fado:password@localhost:5454/fado_db \
  --caddy-admin-url http://localhost:2019 \
  --lb-port 6000
