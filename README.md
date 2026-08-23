# A scalable file processing system built in Go


## Architecture

API Node:
- Stateless
- Workflow:
    - Serves app API and static assets
    - Upload input files to Storage
    - Accept POST /task request from client
    - Generate Task UUID
    - Store job metadata in Redis Hash
    - Push task id to Queue (Redis List)
    - Get task metadata through a GET /task endpoint
    - Download processed files

Task Metadata Storage (Redis Hash):
- Stores file names, instructions, status, creation time and error/output

Task Queue (Redis List):
- Stores task IDs

Worker Node:
- Stateless
- Workflow:
    - Block and listen on Redis List
    - Accept a single task at a time
    - Fetch task metadata from Redis Hash, files from Storage
    - Perform operations on files
    - Store output files to Storage, update task status on redis Hash
    - Repeat


## API Endpoints

- GET / : Get frontend files
- POST /upload : Upload a file for processing
- POST /task : Enqueue a new file processing task
- GET /task/{task_id} : Fetch the details of a task
- GET /download/{filename} : Download a processed output file


## Build & Run using Docker & Docker Compose

Move to root project directory

Build the api and worker images:
```bash
docker build -f api/Dockerfile -t rajsekhar05/file-processor-api:latest .

docker build -f worker/Dockerfile -t rajsekhar05/file-processor-worker:latest .
```

Push the images to Docker Hub:
```bash
docker push rajsekhar05/file-processor-api:latest

docker push rajsekhar05/file-processor-worker:latest
```

OR use script for building & pushing images:
```bash
chmod +x scripts/*
scripts/build-deploy.sh
```

Pull images and run the containers using Docker Compose:
```bash
docker compose up
```


## Running the application

Clone the repository

Run the API Server:
```bash
cd api
go run .
```

Run the worker server:
```bash
cd worker
go run .
```

Send requests to the api server

Test api server:
```bash
cd api
go test -v
```

Test worker server:
```bash
cd worker
go test -v
```