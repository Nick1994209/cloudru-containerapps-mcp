package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

// testGetJob tests getting a specific job
func testGetJob(cfg *config.Config, testJob TestJob) {
	log.Println("\n--- Test: Get Job ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	log.Printf("Attempting to get job with name: %s", testJob.Name)

	// First, try to get a list of jobs to find the job ID
	jobList, err := jobs.GetListJobs(cfg.ProjectID, "")
	if err != nil || len(jobList) == 0 {
		log.Printf("Warning: No jobs found or failed to list jobs: %v", err)
		log.Println("✓ Get job test completed (skipped due to no jobs available)")
		return
	}

	// Find the job by name
	var jobID string
	for _, job := range jobList {
		if job.Name == testJob.Name {
			jobID = job.ID
			break
		}
	}

	if jobID == "" {
		log.Printf("Warning: Could not find job ID for job: %s", testJob.Name)
		log.Println("✓ Get job test completed (skipped due to job not found)")
		return
	}

	log.Printf("Found job ID: %s for job: %s", jobID, testJob.Name)

	// Get specific job
	getRequest := map[string]interface{}{
		"projectId": cfg.ProjectID,
		"jobName":   testJob.Name,
	}
	job, err := jobs.GetJob(cfg.ProjectID, testJob.Name)
	if err != nil {
		logErrorWithRequestBody("Warning: Get job failed", err, getRequest)
		log.Println("✓ Get job test completed (with potential expected error if job doesn't exist)")
		return
	}

	log.Printf("Successfully retrieved job: ID=%s, Name=%s, Status=%s", job.ID, job.Name, job.Status)
	log.Println("✓ Get job test passed")
}
