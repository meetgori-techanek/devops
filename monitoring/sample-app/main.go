package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var logLevels = []string{"INFO", "DEBUG", "WARN", "ERROR", "FATAL"}

var logMessages = []string{
	"User authenticated successfully",
	"Database connection established",
	"Cache miss for key: session_123",
	"Retrying failed request to service A",
	"Received malformed request",
	"Configuration loaded from env",
	"Memory usage exceeds threshold",
	"Invalid login attempt detected",
	"Scheduled job started",
	"Service dependency timeout",
	"File not found: config.yaml",
	"Transaction committed successfully",
	"Unexpected EOF while reading body",
	"Metrics pushed to Prometheus",
	"Cluster leader election triggered",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go cpuBurner()
	go randomLogGenerator()

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is Up and running")
	})

	port := "8080"
	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func cpuBurner() {
	for {
		start := time.Now()
		// Busy loop to consume ~60% CPU
		for time.Since(start) < 60*time.Millisecond {
			_ = rand.Float64() * rand.Float64()
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func randomLogGenerator() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	for {
		level := logLevels[rand.Intn(len(logLevels))]
		msg := logMessages[rand.Intn(len(logMessages))]
		logger.Printf("[%s] %s", level, msg)

		// Sleep random time between 300ms - 1500ms
		time.Sleep(time.Duration(300+rand.Intn(1200)) * time.Millisecond)
	}
}
