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

Requirements: Git, Docker, Docker-Compose, Go 1.25+

1. Clone the repository
```bash
git clone https://github.com/rajsekharde/file-processor.git
cd file-processor
```

2. Build the api and worker images:
```bash
make build-all
```

3. Run the application using Docker Compose:
```bash
make run-app
```

Open localhost:8000 in a browser and perform file operations

4. Stop the application and remove volumes/containers:
```bash
make stop-app
```


## Run the application without building binaries

Clone the repository

Run a standalone Redis container:
```bash
make run-redis

# Stop the container
make stop-redis
```

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

Open localhost:8000 in a browser and perform file operations

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