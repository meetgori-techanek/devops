package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var logChannel = make(chan string, 100) // Shared log channel

func main() {
	go cpuBurner()
	go logGenerator()

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is Up and running")
	})

	http.HandleFunc("/", logStreamHandler)

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
		for time.Since(start) < 60*time.Millisecond {
			_ = rand.Float64() * rand.Float64()
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func logGenerator() {
	logFile := os.Stdout
	logger := log.New(logFile, "", log.LstdFlags)

	levels := []string{"INFO", "WARN", "ERROR"}
	messages := []string{
		"Starting process",
		"Connecting to database",
		"User login failed",
		"Cache miss",
		"Configuration loaded",
		"Retrying request",
	}

	for {
		level := levels[rand.Intn(len(levels))]
		msg := messages[rand.Intn(len(messages))]
		logLine := fmt.Sprintf("[%s] %s", level, msg)

		// Write to stdout
		logger.Println(logLine)

		// Send to log stream
		select {
		case logChannel <- logLine:
		default:
			// Drop log if channel is full
		}

		time.Sleep(time.Second)
	}
}

func logStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Enable real-time flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// Stream logs to the client
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case logLine := <-logChannel:
			fmt.Fprintln(w, logLine)
			flusher.Flush()
		}
	}
}
