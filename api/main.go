package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// root := http.HandlerFunc(HandleRoot)
	root := http.FileServer(http.Dir("../frontend"))
	task := http.HandlerFunc(HandleGetTask)
	postTask := http.HandlerFunc(HandlePostTask)
	uploadFile := http.HandlerFunc(HandleUpload)

	mux.Handle("GET /", LoggerMiddleware(root))
	mux.Handle("GET /task", LoggerMiddleware(task))
	mux.Handle("POST /task", LoggerMiddleware(postTask))
	mux.Handle("POST /upload", LoggerMiddleware(uploadFile))

	log.Printf("API Server running on port 8000...\n\n")
	http.ListenAndServe(":8000", mux)
}
