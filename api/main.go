package main

import (
	"log"
	"net/http"

	"github.com/rajsekharde/file-processor/shared"
)

func main() {
	mux := http.NewServeMux()

	root := http.FileServer(http.Dir("../frontend"))
	task := http.HandlerFunc(HandleGetTask)
	postTask := http.HandlerFunc(HandlePostTask)
	uploadFile := http.HandlerFunc(HandleUpload)

	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("GET /task", shared.LoggerMiddleware(task))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))
	mux.Handle("POST /upload", shared.LoggerMiddleware(uploadFile))

	log.Printf("API Server running on port 8000...\n\n")
	http.ListenAndServe(":8000", mux)
}
