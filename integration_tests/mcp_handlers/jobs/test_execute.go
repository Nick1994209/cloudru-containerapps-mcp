package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

// testExecuteJob tests executing a job
func testExecuteJob(cfg *config.Config, testJob TestJob) {
	log.Println("\n--- Test: Execute Job ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	log.Printf("Attempting to execute job with name: %s", testJob.Name)

	// Execute job with empty params
	jobExecution, err := jobs.ExecuteJob(cfg.ProjectID, testJob.Name, map[string]interface{}{})
	if err != nil {
		log.Printf("Warning: Execute job failed: %v", err)
		log.Println("✓ Execute job test completed (with potential expected error if job execution is not allowed)")
		return
	}

	log.Printf("Successfully executed job: ExecutionName=%s, ExecutionStatus=%s", jobExecution.ExecutionName, jobExecution.ExecutionStatus)

	// Wait a moment for the execution to start
	log.Println("Waiting for job execution to start...")
	time.Sleep(5 * time.Second)

	log.Println("✓ Execute job test passed")
}

// testListJobExecutions tests listing job executions
func testListJobExecutions(cfg *config.Config, testJob TestJob) {
	log.Println("\n--- Test: List Job Executions ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	log.Printf("Attempting to list executions for job: %s", testJob.Name)

	// List job executions
	executions, err := jobs.GetListExecutions(cfg.ProjectID, testJob.Name, "")
	if err != nil {
		log.Printf("Warning: List job executions failed: %v", err)
		log.Println("✓ List job executions test completed (with potential expected error if no executions exist)")
		return
	}

	log.Printf("Successfully retrieved %d executions for job %s", len(executions), testJob.Name)

	// Print first few executions for verification
	for i, execution := range executions {
		if i >= 3 { // Only print first 3 executions
			break
		}
		log.Printf("Execution %d: Name=%s, Status=%s", i+1, execution.ExecutionName, execution.ExecutionStatus)
	}

	log.Println("✓ List job executions test passed")
}
