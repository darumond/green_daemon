package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func clientMain() {
	// Point this to your local server (Safe)
	// OR use https://httpbin.org/post (Safe for testing)
	targetURL := "http://localhost:8080"
	interval := 2 * time.Second

	payload := []byte("Hello, Server! Same size every time.")
	client := &http.Client{Timeout: 5 * time.Second}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		req, _ := http.NewRequest("POST", targetURL, bytes.NewReader(payload))

		// Adding a User-Agent makes the request look more legitimate
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0.0.0")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("[%s] Sent %d bytes. Status: %s\n", time.Now().Format("15:04:05"), len(payload), resp.Status)
		resp.Body.Close()
	}
}
