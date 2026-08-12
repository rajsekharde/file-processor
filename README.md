# A scalable file processing system built in Go

# Running the application

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

Run the containers using Docker Compose:
```bash
docker compose up
```


## Architecture

API Node:
- Stateless
- Connects to DB, Task Queue and File Storage using interfaces
- Workflow:
    - Accept request from client
    - Upload files to Storage
    - Store job metadata in DB
    - Push task data to Queue

Worker Node:
- Stateless
- Connects to DB, Task Queue and File Storage using interfaces
- Workflow:
    - Block and listen on Queue
    - Accept a single task at a time
    - Fetch task metadata from DB, files from Storage
    - Process files
    - Store files to DB, update job status on DB
    - Repeat

Local Operations:
- DB: PostgreSQL running on Docker
- Task Queue: Redis running on Docker
- File Storage: Shared Docker volume

## Tasks

Learn Go concepts:
- File I/O
- Logging middleware
- Reading and parsing request body
- Writing JSON data to response
- Calling external API and parsing its data