# A scalable file processing system built in Go

## Architecture

API Node:
- Stateless, horizontally scalable
- Connects to DB, Task Queue and File Storage using interfaces
- Workflow:
    - Accept request from client
    - Upload files to Storage
    - Store job metadata in DB
    - Push task data to Queue

Worker Node:
- Stateless, horizontally scalable
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