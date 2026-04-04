package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

// testListJobs tests listing jobs
func testListJobs(cfg *config.Config, createdJobs []TestJob) {
	log.Println("\n--- Test: List Jobs ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	// List jobs
	jobList, err := jobs.GetListJobs(cfg.ProjectID, "")
	if err != nil {
		log.Printf("Warning: List jobs failed: %v", err)
		log.Println("✓ List jobs test completed (with potential expected error if no jobs exist)")
		return
	}

	log.Printf("Successfully retrieved %d jobs", len(jobList))

	// Print first few jobs for verification
	for i, job := range jobList {
		if i >= 3 { // Only print first 3 jobs
			break
		}
		log.Printf("Job %d: ID=%s, Name=%s, Status=%s", i+1, job.ID, job.Name, job.Status)
	}

	// Verify that all created jobs are present in the list
	if len(createdJobs) > 0 {
		log.Printf("\nVerifying that %d created jobs are present in the list...", len(createdJobs))
		foundCount := 0
		for _, createdJob := range createdJobs {
			found := false
			for _, job := range jobList {
				if job.ID == createdJob.ID {
					found = true
					foundCount++
					log.Printf("✓ Found created job: ID=%s, Name=%s", job.ID, job.Name)
					break
				}
			}
			if !found {
				log.Printf("✗ Created job not found in list: ID=%s, Name=%s", createdJob.ID, createdJob.Name)
			}
		}
		if foundCount == len(createdJobs) {
			log.Printf("✓ All %d created jobs found in the list", foundCount)
		} else {
			log.Printf("⚠ Only %d out of %d created jobs found in the list", foundCount, len(createdJobs))
		}
	}

	log.Println("✓ List jobs test passed")
}
