#!/bin/bash

# Run a containerized PostgreSQL data from the local .sql files.

cd $(dirname $0)
cd ../database
ENTRYPOINT=`pwd`

docker run \
  --name standalone-fado-db \
  -p 5454:5432 \
  -v ${ENTRYPOINT}:/docker-entrypoint-initdb.d \
  -e POSTGRES_USER=fado \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=fado_db \
  --rm \
  postgres:13
