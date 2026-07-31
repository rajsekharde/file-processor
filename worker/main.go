package main

import (
	"log"
	"net/http"

	"github.com/rajsekharde/file-processor/shared"
)

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(handleRoot)
	getTask := http.HandlerFunc(handleTask)
	postTask := http.HandlerFunc(handlePostTask)
	
	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("GET /task", shared.LoggerMiddleware(getTask))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))

	log.Printf("Worker Server running on port 8001...\n\n")
	http.ListenAndServe(":8001", mux)
}