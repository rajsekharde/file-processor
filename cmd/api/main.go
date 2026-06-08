package main

import (
	"fmt"
	"io"
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8001/task")
	var msg string

	if err != nil {
		log.Println(err.Error())
		msg = "Could not fetch API"
	} else {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
		}
		msg = string(responseData)
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, msg)
}

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(handleRoot)
	task := http.HandlerFunc(handleGetTask)

	mux.Handle("GET /", loggerMiddleware(root))
	mux.Handle("GET /task", loggerMiddleware(task))

	log.Println("API Server running on port 8000...")
	http.ListenAndServe(":8000", mux)
}
