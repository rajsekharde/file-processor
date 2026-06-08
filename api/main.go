package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleRoot)
	log.Println("Server running on port 8000")
	http.ListenAndServe(":8000", mux)
}