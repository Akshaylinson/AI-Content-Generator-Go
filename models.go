package main

import (
	"time"
)

type Job struct {
	ID        int       `json:"id"`
	Topic     string    `json:"topic"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateJobRequest struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
