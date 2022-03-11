#!/bin/bash

# Connect to the database spun up by 'run-db.sh'.

PGPASSWORD="password" psql -h localhost -p 5454 -d fado_db -U fado