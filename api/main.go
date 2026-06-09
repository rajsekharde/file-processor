package main

import (
	"log"
	"net/http"
	"time"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Pass control to the next handler in the chain
		next.ServeHTTP(w, r)
		
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(HandleRoot)
	task := http.HandlerFunc(HandleGetTask)

	mux.Handle("GET /", loggerMiddleware(root))
	mux.Handle("GET /task", loggerMiddleware(task))

	log.Println("API Server running on port 8000...")
	http.ListenAndServe(":8000", mux)
}
