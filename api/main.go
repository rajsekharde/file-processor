package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(HandleRoot)
	task := http.HandlerFunc(HandleGetTask)
	postTask := http.HandlerFunc(HandlePostTask)

	mux.Handle("GET /", LoggerMiddleware(root))
	mux.Handle("GET /task", LoggerMiddleware(task))
	mux.Handle("POST /task", LoggerMiddleware(postTask))

	log.Println("API Server running on port 8000...")
	http.ListenAndServe(":8000", mux)
}
