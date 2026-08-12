#!/bin/bash

# Build and Deploy api and worker images

echo "Building images"

docker build -f api/Dockerfile -t rajsekhar05/file-processor-api:latest .

docker build -f worker/Dockerfile -t rajsekhar05/file-processor-worker:latest .

echo "Images built"

echo "Deploying images"

docker push rajsekhar05/file-processor-api:latest

docker push rajsekhar05/file-processor-worker:latest

echo "Images deployed"