package main

import (
	"fmt"
	"net/http"
)

func serverMain() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received %s request from %s\n", r.Method, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Acknowledged"))
	})

	fmt.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
