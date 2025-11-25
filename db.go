package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "../db/content.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic TEXT NOT NULL,
		type TEXT DEFAULT 'blog',
		status TEXT NOT NULL DEFAULT 'pending',
		output TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_status ON jobs(status);
	CREATE INDEX IF NOT EXISTS idx_created_at ON jobs(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

func createJob(topic, jobType string) (*Job, error) {
	query := `INSERT INTO jobs (topic, type, status, created_at, updated_at) VALUES (?, ?, 'pending', ?, ?)`
	now := time.Now()
	
	if jobType == "" {
		jobType = "blog"
	}
	
	result, err := db.Exec(query, topic, jobType, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get job ID: %v", err)
	}

	return &Job{
		ID:        int(id),
		Topic:     topic,
		Type:      jobType,
		Status:    "pending",
		Output:    "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func getJobs() ([]Job, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	query := `SELECT id, topic, COALESCE(type, 'blog'), status, output, created_at, updated_at FROM jobs ORDER BY created_at DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %v", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Topic, &job.Type, &job.Status, &job.Output, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %v", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func getJobByID(id int) (*Job, error) {
	query := `SELECT id, topic, COALESCE(type, 'blog'), status, output, created_at, updated_at FROM jobs WHERE id = ?`
	
	var job Job
	err := db.QueryRow(query, id).Scan(&job.ID, &job.Topic, &job.Type, &job.Status, &job.Output, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %v", err)
	}

	return &job, nil
}

func deleteJob(id int) error {
	query := `DELETE FROM jobs WHERE id = ?`
	
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func updateJobStatus(jobID int, status, output string) error {
	query := `UPDATE jobs SET status = ?, output = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	
	result, err := db.Exec(query, status, output, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}
