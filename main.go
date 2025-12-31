package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FibonacciResponse struct {
	OrderNumber  int    `json:"order_number"`
	FibonacciNum uint64 `json:"fibonacci_number"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message,omitempty"`
}

// Time: O(2^n) - extremely slow for n > 40
func recursiveFibonacci(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	return recursiveFibonacci(n-1) + recursiveFibonacci(n-2)
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orderStr := r.PathValue("order_number")

	order, err := strconv.Atoi(orderStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if order < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if order > 93 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	fibNum := recursiveFibonacci(order)
	elapsed := time.Since(startTime)

	response := FibonacciResponse{
		OrderNumber:  order,
		FibonacciNum: fibNum,
		StatusCode:   http.StatusOK,
		Message:      fmt.Sprintf("Computed in %v", elapsed),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Printf("Fibonacci(%d) = %d (computed in %v)", order, fibNum, elapsed)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("Completed %s %s - %d in %v",
			r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	// Create a new router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /api/v1/fibonacci/{order_number}", fibonacciHandler)

	// Apply middleware
	handler := loggingMiddleware(mux)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	log.Printf("Starting Fibonacci API server on port 8080")
	log.Printf("Try these examples:")
	log.Printf("  http://localhost:8080/api/v1/fibonacci/3")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
