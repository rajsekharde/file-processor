package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8001/task")

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not fetch API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not read data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(responseData))
}
