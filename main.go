package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var worker *ContentWorker

func main() {
	// Create db directory if it doesn't exist
	if err := os.MkdirAll("../db", 0755); err != nil {
		log.Fatalf("Failed to create db directory: %v", err)
	}

	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Start content worker
	worker = NewContentWorker()
	worker.Start()
	defer worker.Stop()

	// Setup routes
	r := mux.NewRouter()
	
	// API routes - direct paths
	r.HandleFunc("/api/jobs", getJobsHandler).Methods("GET")
	r.HandleFunc("/api/jobs", createJobHandler).Methods("POST")
	r.HandleFunc("/api/job/{id}", getJobHandler).Methods("GET")
	r.HandleFunc("/api/job/{id}", deleteJobHandler).Methods("DELETE")
	r.HandleFunc("/api/process", processJobsHandler).Methods("POST")
	r.HandleFunc("/api/model-status", modelStatusHandler).Methods("GET")
	
	// Dashboard route
	r.HandleFunc("/", dashboardHandler).Methods("GET")
	
	// Add CORS middleware
	r.Use(corsMiddleware)

	log.Println("Server starting on :8080")
	log.Println("Dashboard: http://localhost:8080")
	log.Println("API: http://localhost:8080/api/jobs")
	
	// Check for model
	if modelPath := findModel(); modelPath != "" {
		log.Printf("Local model detected: %s", modelPath)
	} else {
		log.Println("No local model found - using enhanced fallback generation")
	}
	
	// Log all routes for debugging
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s %v", pathTemplate, methods)
		return nil
	})
	
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
