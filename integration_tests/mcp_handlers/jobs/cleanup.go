package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

// cleanupTestJobs deletes all jobs that start with "test-job-"
func cleanupTestJobs(cfg *config.Config) {
	log.Println("\n--- Cleanup: Deleting Test Jobs ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	// List all jobs
	jobList, err := jobs.GetListJobs(cfg.ProjectID, "")
	if err != nil {
		log.Printf("Warning: Failed to list jobs for cleanup: %v", err)
		log.Println("✓ Cleanup completed (skipped due to list failure)")
		return
	}

	if len(jobList) == 0 {
		log.Println("No jobs found for cleanup")
		log.Println("✓ Cleanup completed")
		return
	}

	// Find and delete all test jobs
	deletedCount := 0
	for _, job := range jobList {
		if strings.HasPrefix(job.Name, "test-job-") {
			log.Printf("Deleting test job: %s (ID: %s)", job.Name, job.ID)
			deleteRequest := map[string]interface{}{
				"projectId": cfg.ProjectID,
				"jobName":   job.Name,
			}
			_, err := jobs.DeleteJob(cfg.ProjectID, job.Name)
			if err != nil {
				logErrorWithRequestBody(fmt.Sprintf("Warning: Failed to delete job %s", job.Name), err, deleteRequest)
			} else {
				deletedCount++
				log.Printf("Successfully deleted test job: %s", job.Name)
			}
		}
	}

	if deletedCount > 0 {
		log.Printf("✓ Cleanup completed: Deleted %d test job(s)", deletedCount)
	} else {
		log.Println("✓ Cleanup completed: No test jobs found to delete")
	}
}
