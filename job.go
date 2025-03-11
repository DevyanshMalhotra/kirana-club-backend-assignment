package main

import (
	"sync"
)

type JobStatus string

const (
	StatusOngoing   JobStatus = "ongoing"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID     int         `json:"job_id"`
	Status JobStatus   `json:"status"`
	Visits []Visit     `json:"visits"`
	Errors []JobError  `json:"error,omitempty"`
}

type Visit struct {
	StoreID   string        `json:"store_id"`
	VisitTime string        `json:"visit_time"`
	ImageURLs []string      `json:"image_url"`
	Images    []ImageResult `json:"images,omitempty"`
}

type ImageResult struct {
	URL       string `json:"url"`
	Perimeter int    `json:"perimeter,omitempty"`
	Error     string `json:"error,omitempty"`
}

type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

var (
	jobIDCounter = 1
	jobs         = make(map[int]*Job)
	jobsMutex    sync.Mutex
)

func createJob(visits []Visit) *Job { // creates a new job and stores it in memory
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	job := &Job{
		ID:     jobIDCounter,
		Status: StatusOngoing,
		Visits: visits,
	}
	jobs[jobIDCounter] = job
	jobIDCounter++
	return job
}


func getJob(jobID int) (*Job, bool) { // retrieves a job by its id
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	job, exists := jobs[jobID]
	return job, exists
}

func updateJob(job *Job) { // updates the job stored in memory
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	jobs[job.ID] = job
}
