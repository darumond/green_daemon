package main

import (
	"green-daemon/procinfo"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "green_daemon_http_requests_total",
		Help: "Total HTTP requests to root endpoint",
	},
)

func main() {
	// Register metric
	prometheus.MustRegister(httpRequestsTotal)

	// Start your polling loop
	go procinfo.PollLoop(time.Second * 2)

	// Add a basic endpoint to generate traffic
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.Inc()
		w.Write([]byte("Green daemon running"))
	})

	http.Handle("/metrics", promhttp.Handler())

	println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
