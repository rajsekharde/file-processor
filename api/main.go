package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, msg)
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /task", handleGetTask)

	log.Println("API Server running on port 8000...")
	http.ListenAndServe(":8000", mux)
}
