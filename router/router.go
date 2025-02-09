package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TimiBolu/lema-ai-users-service/config"
	"github.com/TimiBolu/lema-ai-users-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// ANSI color codes for terminal output
const (
	green  = "\033[32m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

// ResponseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getAPIDocs(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("docs/api.md")
	if err != nil {
		http.Error(w, "Failed to load documentation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Logging middleware to track API calls and response times with color
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Wrap the ResponseWriter to capture the status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(rw, r)

		// Determine color for method
		var methodColor string
		switch r.Method {
		case "GET":
			methodColor = green
		case "POST":
			methodColor = blue
		case "DELETE":
			methodColor = yellow
		default:
			methodColor = reset
		}

		// Determine color for status code
		var statusColor string
		switch {
		case rw.statusCode >= 200 && rw.statusCode < 300:
			statusColor = green // ✅ Success
		case rw.statusCode >= 400 && rw.statusCode < 500:
			statusColor = yellow // ⚠️ Client Error
		case rw.statusCode >= 500:
			statusColor = red // ❌ Server Error
		default:
			statusColor = reset
		}

		// Log with colors
		duration := time.Since(startTime)
		log.Printf("📡 %s%s%s %s | Status: %s%d%s | ⏱️ %v",
			methodColor, r.Method, reset,
			r.URL.Path,
			statusColor, rw.statusCode, reset,
			duration,
		)
	})
}

func Setup() {
	// Initialize the router
	r := mux.NewRouter()

	// Add logging middleware
	r.Use(loggingMiddleware)

	// User endpoints
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/count", handlers.GetUsersCount).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserByID).Methods("GET")

	// Post endpoints
	r.HandleFunc("/posts", handlers.GetPostsByUser).Methods("GET")
	r.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	// API Documentation endpoints
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/index.html")
	})
	r.HandleFunc("/docs/raw", getAPIDocs) // Serve raw Markdown

	frontendApps := config.EnvConfig.FRONTEND_APPS
	allowedOrigins := strings.Split(frontendApps, ",")

	// Apply CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	// Server configuration
	port := config.EnvConfig.PORT
	baseURL := config.EnvConfig.SERVER_BASE_URL
	log.Printf("🚀 Server is up and running on %s:%s", baseURL, port)
	log.Printf("📄 API Documentation available at %s:%s/docs", baseURL, port)

	// Start the server with CORS
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler)
	if err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
