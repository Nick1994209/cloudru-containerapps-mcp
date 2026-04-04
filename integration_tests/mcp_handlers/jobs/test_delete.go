package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

// testDeleteJob tests deleting a job
func testDeleteJob(cfg *config.Config, testJobs []TestJob) {
	log.Println("\n--- Test: Delete Job ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	if len(testJobs) == 0 {
		log.Println("No jobs to delete, skipping test")
		log.Println("✓ Delete job test completed (skipped)")
		return
	}

	// Delete each test job
	for _, testJob := range testJobs {
		log.Printf("\nDeleting job: %s", testJob.Name)

		// Get the job ID by listing jobs
		jobList, err := jobs.GetListJobs(cfg.ProjectID, "")
		if err != nil || len(jobList) == 0 {
			log.Printf("Warning: Failed to list jobs to find job ID: %v", err)
			continue
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
			continue
		}

		log.Printf("Found job ID: %s for job: %s", jobID, testJob.Name)

		// Test 1: Delete the job
		log.Printf("\nTest 1: Delete job %s", testJob.Name)
		deleteRequest := map[string]interface{}{
			"projectId": cfg.ProjectID,
			"jobName":   testJob.Name,
		}
		operation, err := jobs.DeleteJob(cfg.ProjectID, testJob.Name)
		if err != nil {
			logErrorWithRequestBody("Warning: Delete job failed", err, deleteRequest)
		} else {
			if validateOperation(operation) {
				log.Printf("Successfully deleted job: %s (ID: %s), Operation ID: %s", testJob.Name, jobID, operation.ResourceID)
			} else {
				log.Printf("Warning: Operation validation failed for delete operation on job: %s", testJob.Name)
			}
		}

		// Test 2: Verify job is deleted by trying to get it
		log.Printf("\nTest 2: Verify job %s is deleted", testJob.Name)
		getRequest := map[string]interface{}{
			"projectId": cfg.ProjectID,
			"jobName":   testJob.Name,
		}
		_, err = jobs.GetJob(cfg.ProjectID, testJob.Name)
		if err != nil {
			logErrorWithRequestBody("Job successfully deleted (GetJob returned error as expected)", err, getRequest)
		} else {
			log.Printf("Warning: Job still exists after deletion")
		}

		// Test 3: Try to delete the same job again (should fail)
		log.Printf("\nTest 3: Try to delete already deleted job %s", testJob.Name)
		deleteRequest2 := map[string]interface{}{
			"projectId": cfg.ProjectID,
			"jobName":   testJob.Name,
		}
		_, err = jobs.DeleteJob(cfg.ProjectID, testJob.Name)
		if err != nil {
			logErrorWithRequestBody("Expected error when deleting already deleted job", err, deleteRequest2)
		} else {
			log.Printf("Warning: Delete operation succeeded for already deleted job")
		}
	}

	log.Println("✓ Delete job test completed")
}
