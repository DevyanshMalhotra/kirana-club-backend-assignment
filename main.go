package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {
	err := LoadStoreMaster("StoreMasterAssignment.csv")
	if err != nil {
		log.Fatalf("Failed to load store master: %v", err)
	}

	http.HandleFunc("/api/submit/", submitJobHandler)
	http.HandleFunc("/api/status", getJobStatusHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitJobHandler(w http.ResponseWriter, r *http.Request) { //accepts job submissions
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Count  int     `json:"count"`
		Visits []Visit `json:"visits"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Count != len(payload.Visits) {
		http.Error(w, `{"error": "count does not match visits"}`, http.StatusBadRequest)
		return
	}

	job := createJob(payload.Visits)
	go processJob(job)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"job_id": job.ID})
}

func getJobStatusHandler(w http.ResponseWriter, r *http.Request) { //returns job status 
	jobidStr := r.URL.Query().Get("jobid")
	if jobidStr == "" {
		http.Error(w, "jobid parameter missing", http.StatusBadRequest)
		return
	}
	jobid, err := strconv.Atoi(jobidStr)
	if err != nil {
		http.Error(w, "invalid jobid", http.StatusBadRequest)
		return
	}

	job, exists := getJob(jobid)
	if !exists {
		http.Error(w, "job not found", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"job_id": job.ID,
		"status": job.Status,
	}
	if job.Status == StatusFailed {
		response["error"] = job.Errors
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
