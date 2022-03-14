#!/bin/bash

# Run the current frontend client.
# Make sure to 'npm install' beforehand.

cd $(dirname $0)
cd ../client

npm run start
