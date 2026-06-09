package main

import (
	"net/http"
	"time"
	"log"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Pass control to the next handler in the chain
		next.ServeHTTP(w, r)
		
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}