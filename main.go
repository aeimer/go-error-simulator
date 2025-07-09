package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Configuration
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		_, err := strconv.Atoi(p)
		if err != nil {
			slog.Error("Invalid PORT environment variable, using default port 8080", "error", err)
		} else {
			port = p
		}
	}

	// HTTP handler
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/simulate", simulateHandler)

	// Starting server
	slog.Info("Starting server", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// HTTP handler functions
func rootHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Serving root request", "method", r.Method, "path", r.URL.Path)

	text := `Welcome to the go-error-simulator!
Use the /simulate endpoint to trigger custom responses and logs.

Examples:
  GET /simulate?status=200
  GET /simulate?status=404&stdout_msg=File_not_found
  GET /simulate?status=500&stderr_msg=Internal_server_error&stdout_msg=Processing_failed
  GET /simulate?status=401&stderr_msg=Auth_failed
  GET /simulate?status=401&latency=300-900
`
	_, _ = fmt.Fprint(w, text)
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Serving simulate request", "method", r.Method, "path", r.URL.Path, "raddr", r.RemoteAddr, "query", r.URL.Query())

	latency := latency(r)

	// --- Simulate HTTP Status Code ---
	statusCodeStr := r.URL.Query().Get("status")
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil || statusCode < 100 || statusCode >= 600 {
		slog.Error("Invalid or missing 'status' parameter. Defaulting to 400.", "status_code", statusCodeStr, "error", err)
		statusCode = http.StatusBadRequest
	}

	// --- Write to stdout ---
	stdoutMsg := r.URL.Query().Get("stdout_msg")
	if stdoutMsg != "" {
		fmt.Println("STDOUT_TRIGGER:", stdoutMsg)
	}

	// --- Write to stderr ---
	stderrMsg := r.URL.Query().Get("stderr_msg")
	if stderrMsg != "" {
		_, _ = fmt.Fprintln(os.Stderr, "STDERR_TRIGGER:", stderrMsg)
	}

	// --- Set response headers ---
	w.WriteHeader(statusCode)
	responseMessage := fmt.Sprintf(
		"Simulated HTTP %d\nSTDOUT triggered: %t\nSTDERR triggered: %t\nLatency: %dms",
		statusCode, stdoutMsg != "", stderrMsg != "", latency,
	)
	_, _ = fmt.Fprintln(w, responseMessage)

	slog.Info("Responded simulate", "status_code", statusCode, "path", r.URL.Path, "latency", latency)
}

// Helper

func latency(r *http.Request) int {
	// --- Simulate Latency ---
	latencyMinMs, latencyMaxMs := 0, 0
	latencyStr := r.URL.Query().Get("latency")
	latencySplit := strings.Split(latencyStr, "-")

	if len(latencySplit) > 2 {
		slog.Error("Invalid 'latency' parameter. Expected format is 'min-max'. Defaulting to 0ms.", "latency", latencyStr)
	}

	// Latency is exactly specified as "min-max"
	if len(latencySplit) == 2 {
		var err error
		latencyMinMs, err = strconv.Atoi(latencySplit[0])
		if err != nil {
			slog.Error("Invalid 'latency' parameter. Defaulting to 0ms.", "latency", latencyStr, "error", err)
			latencyMinMs = 0
		}

		latencyMaxMs, err = strconv.Atoi(latencySplit[1])
		if err != nil {
			slog.Error("Invalid 'latency' parameter. Defaulting to 0ms.", "latency", latencyStr, "error", err)
			latencyMaxMs = 0
		}
	}

	if len(latencySplit) == 1 {
		if latencyStr[:1] == "-" {
			// Format is like "-100" means 0 to 100ms latency
			var err error
			latencyMaxMs, err = strconv.Atoi(latencyStr[1:])
			if err != nil {
				slog.Error("Invalid 'latency' parameter. Defaulting to 0ms.", "latency", latencyStr, "error", err)
				latencyMaxMs = 0
			}
		} else {
			//	Format is expected to be like "100" means exactly 100ms latency
			var err error
			latencyMinMs, err = strconv.Atoi(latencyStr)
			if err != nil {
				slog.Error("Invalid 'latency' parameter. Defaulting to 0ms.", "latency", latencyStr, "error", err)
			}
			latencyMaxMs = latencyMinMs
		}
	}

	// Validate latency range
	if 0 > latencyMinMs || latencyMinMs > latencyMaxMs {
		slog.Error("'latency' parameter min cannot be greater than max and negative values are not allowed. Defaulting to 0ms.\n", "latency_min", latencyMinMs, "latency_max", latencyMaxMs)
		latencyMinMs = 0
		latencyMaxMs = 0
	}

	// Default value is min as we expect latency to be defined exactly
	latency := latencyMinMs

	// If the latency is specified as a range, calculate a random latency within that range
	if latencyMaxMs != latencyMinMs {
		latency = latencyMinMs + rand.Intn(latencyMaxMs-latencyMinMs+1)
	}
	time.Sleep(time.Duration(latency) * time.Millisecond)
	return latency
}
