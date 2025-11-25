package main

import (
	"database/sql"
	"log"
	"time"
)

type ContentWorker struct {
	generator *LLMGenerator
	running   bool
}

func NewContentWorker() *ContentWorker {
	return &ContentWorker{
		generator: NewLLMGenerator(),
		running:   true,
	}
}

func (w *ContentWorker) Start() {
	log.Println("Content worker started")
	go w.run()
}

func (w *ContentWorker) Stop() {
	w.running = false
	log.Println("Content worker stopped")
}

func (w *ContentWorker) run() {
	log.Println("Worker loop started")
	for w.running {
		job := w.getPendingJob()
		if job != nil {
			log.Printf("Found pending job: %d", job.ID)
			w.processJob(job)
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	log.Println("Worker loop stopped")
}

func (w *ContentWorker) ProcessAllPending() {
	jobs := w.getAllPendingJobs()
	log.Printf("Processing %d pending jobs", len(jobs))
	
	for _, job := range jobs {
		w.processJob(&job)
	}
}

func (w *ContentWorker) getPendingJob() *Job {
	query := `SELECT id, topic, COALESCE(type, 'blog') FROM jobs WHERE status = 'pending' ORDER BY created_at ASC LIMIT 1`
	
	var job Job
	err := db.QueryRow(query).Scan(&job.ID, &job.Topic, &job.Type)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error getting pending job: %v", err)
		}
		return nil
	}
	
	return &job
}

func (w *ContentWorker) getAllPendingJobs() []Job {
	query := `SELECT id, topic, COALESCE(type, 'blog') FROM jobs WHERE status = 'pending' ORDER BY created_at ASC`
	
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error getting pending jobs: %v", err)
		return nil
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Topic, &job.Type)
		if err != nil {
			log.Printf("Error scanning job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs
}

func (w *ContentWorker) processJob(job *Job) {
	log.Printf("Processing job %d: %s", job.ID, job.Topic)
	
	// Update status to processing
	if err := updateJobStatus(job.ID, "processing", ""); err != nil {
		log.Printf("Failed to update job %d to processing: %v", job.ID, err)
		return
	}
	
	// Generate content
	log.Printf("Generating content for job %d", job.ID)
	content := w.generator.GenerateContent(job.Topic)
	log.Printf("Generated content length: %d characters", len(content))
	
	if content != "" {
		if err := updateJobStatus(job.ID, "completed", content); err != nil {
			log.Printf("Failed to update job %d to completed: %v", job.ID, err)
		} else {
			log.Printf("Job %d completed successfully", job.ID)
		}
	} else {
		updateJobStatus(job.ID, "failed", "Content generation failed")
		log.Printf("Job %d failed - no content generated", job.ID)
	}
}
