package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Topic) == "" {
		writeErrorResponse(w, "Topic is required", http.StatusBadRequest)
		return
	}

	job, err := createJob(req.Topic, req.Type)
	if err != nil {
		log.Printf("Error creating job: %v", err)
		writeErrorResponse(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, job)
}

func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/jobs called")
	
	if db == nil {
		log.Printf("Database not initialized")
		writeErrorResponse(w, "Database not available", http.StatusInternalServerError)
		return
	}
	
	jobs, err := getJobs()
	if err != nil {
		log.Printf("Error getting jobs: %v", err)
		writeErrorResponse(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	log.Printf("Returning %d jobs", len(jobs))
	writeSuccessResponse(w, jobs)
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := getJobByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, "Job not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting job: %v", err)
		writeErrorResponse(w, "Failed to get job", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, job)
}

func deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	err = deleteJob(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, "Job not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting job: %v", err)
		writeErrorResponse(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, map[string]string{"message": "Job deleted successfully"})
}

func processJobsHandler(w http.ResponseWriter, r *http.Request) {
	// Trigger immediate processing of all pending jobs
	go worker.ProcessAllPending()
	
	writeSuccessResponse(w, map[string]string{
		"message": "Processing triggered for all pending jobs",
	})
}

func modelStatusHandler(w http.ResponseWriter, r *http.Request) {
	modelPath := findModel()
	var message string
	
	if modelPath != "" {
		message = fmt.Sprintf("✅ Local model active: %s", filepath.Base(modelPath))
	} else {
		message = "⚠️ No local model found - using enhanced fallback generation"
	}
	
	writeSuccessResponse(w, map[string]string{"message": message})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Dashboard requested")
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>AI Content Automator</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        button { padding: 10px 20px; margin: 10px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        input { padding: 10px; margin: 10px; width: 300px; border: 1px solid #ddd; border-radius: 4px; }
        .job { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🤖 AI Content Automator</h1>
        
        <div>
            <h3>Create New Job</h3>
            <input type="text" id="topic" placeholder="Enter topic..." />
            <button onclick="createJob()">Create Job</button>
            <button onclick="processJobs()">Process All</button>
            <button onclick="loadJobs()">Refresh</button>
        </div>
        
        <div id="jobs"></div>
    </div>
    
    <script>
        async function createJob() {
            const topic = document.getElementById('topic').value;
            if (!topic) return alert('Enter a topic');
            
            try {
                const response = await fetch('/api/jobs', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ topic, type: 'blog' })
                });
                const data = await response.json();
                if (data.success) {
                    document.getElementById('topic').value = '';
                    loadJobs();
                } else {
                    alert('Error: ' + data.error);
                }
            } catch (e) {
                alert('Failed to create job');
            }
        }
        
        async function processJobs() {
            try {
                await fetch('/api/process', { method: 'POST' });
                setTimeout(loadJobs, 2000);
            } catch (e) {
                alert('Failed to process jobs');
            }
        }
        
        async function loadJobs() {
            try {
                const response = await fetch('/api/jobs');
                const data = await response.json();
                if (data.success) {
                    const jobsDiv = document.getElementById('jobs');
                    jobsDiv.innerHTML = '<h3>Jobs (' + data.data.length + ')</h3>';
                    data.data.forEach(job => {
                        jobsDiv.innerHTML += '<div class="job"><strong>#' + job.id + '</strong> - ' + job.topic + '<br><em>Status: ' + job.status + '</em><br>' + (job.output ? job.output.substring(0, 200) + '...' : 'No output yet') + '</div>';
                    });
                }
            } catch (e) {
                document.getElementById('jobs').innerHTML = '<p>Error loading jobs</p>';
            }
        }
        
        // Auto-refresh
        setInterval(loadJobs, 10000);
        loadJobs();
    </script>
</body>
</html>`))
}

func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Success: true,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	json.NewEncoder(w).Encode(response)
}
