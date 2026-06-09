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
		defer resp.Body.Close()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, msg)
}
