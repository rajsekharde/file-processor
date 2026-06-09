package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(HandleRoot)
	task := http.HandlerFunc(HandleGetTask)

	mux.Handle("GET /", LoggerMiddleware(root))
	mux.Handle("GET /task", LoggerMiddleware(task))

	log.Println("API Server running on port 8000...")
	http.ListenAndServe(":8000", mux)
}
