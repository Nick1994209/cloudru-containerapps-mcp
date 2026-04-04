package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// testPatchJob tests patching a job with different parameter combinations
func testPatchJob(cfg *config.Config, testJob TestJob) {
	log.Println("\n--- Test: Patch Job ---")

	jobs := cloudru.NewJobsApplication(cfg)

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	log.Printf("Patching job: %s", testJob.Name)

	// Wait for job to be created
	time.Sleep(5 * time.Second)

	// Test 1: Patch job with new image
	log.Println("\nTest 1: Patch job with new image")
	newImage := "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2"
	patchRequest1 := domain.PatchJobRequest{
		ProjectID: cfg.ProjectID,
		JobName:   testJob.Name,
		JobImage:  &newImage,
	}

	operation, err := jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest1)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new image failed", err, patchRequest1)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new image: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 2: Patch job with new CPU
	log.Println("\nTest 2: Patch job with new CPU")
	newCPU := "0.2"
	patchRequest2 := domain.PatchJobRequest{
		ProjectID: cfg.ProjectID,
		JobName:   testJob.Name,
		JobCPU:    &newCPU,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest2)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new CPU failed", err, patchRequest2)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new CPU: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 3: Patch job with new description
	log.Println("\nTest 3: Patch job with new description")
	newDescription := "Updated job description"
	patchRequest3 := domain.PatchJobRequest{
		ProjectID:      cfg.ProjectID,
		JobName:        testJob.Name,
		JobDescription: &newDescription,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest3)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new description failed", err, patchRequest3)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new description: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 4: Patch job with new environment variables
	log.Println("\nTest 4: Patch job with new environment variables")
	newEnvVars := "NEW_VAR='new_value';ANOTHER_VAR='another_value'"
	patchRequest4 := domain.PatchJobRequest{
		ProjectID:               cfg.ProjectID,
		JobName:                 testJob.Name,
		JobEnvironmentVariables: &newEnvVars,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest4)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new environment variables failed", err, patchRequest4)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new environment variables: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 5: Patch job with new command
	log.Println("\nTest 5: Patch job with new command")
	newCommand := []string{"/bin/sh", "-c"}
	patchRequest5 := domain.PatchJobRequest{
		ProjectID:  cfg.ProjectID,
		JobName:    testJob.Name,
		JobCommand: newCommand,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest5)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new command failed", err, patchRequest5)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new command: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 6: Patch job with new args
	log.Println("\nTest 6: Patch job with new args")
	newArgs := []string{"echo", "Patched job"}
	patchRequest6 := domain.PatchJobRequest{
		ProjectID: cfg.ProjectID,
		JobName:   testJob.Name,
		JobArgs:   newArgs,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest6)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new args failed", err, patchRequest6)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new args: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 7: Patch job with new retry count
	log.Println("\nTest 7: Patch job with new retry count")
	newRetryCount := uint32(10)
	patchRequest7 := domain.PatchJobRequest{
		ProjectID:     cfg.ProjectID,
		JobName:       testJob.Name,
		JobRetryCount: &newRetryCount,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest7)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new retry count failed", err, patchRequest7)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new retry count: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 8: Patch job with new execution timeout
	log.Println("\nTest 8: Patch job with new execution timeout")
	newExecutionTimeout := uint32(1800)
	patchRequest8 := domain.PatchJobRequest{
		ProjectID:           cfg.ProjectID,
		JobName:             testJob.Name,
		JobExecutionTimeout: &newExecutionTimeout,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest8)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new execution timeout failed", err, patchRequest8)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new execution timeout: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 9: Patch job with new privileged setting
	log.Println("\nTest 9: Patch job with new privileged setting")
	newPrivileged := true
	patchRequest9 := domain.PatchJobRequest{
		ProjectID:     cfg.ProjectID,
		JobName:       testJob.Name,
		JobPrivileged: &newPrivileged,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest9)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with new privileged setting failed", err, patchRequest9)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with new privileged setting: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 10: Patch job with multiple parameters
	log.Println("\nTest 10: Patch job with multiple parameters")
	multiImage := "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2"
	multiCPU := "0.3"
	multiDescription := "Multi-patch test"
	multiRequest := domain.PatchJobRequest{
		ProjectID:      cfg.ProjectID,
		JobName:        testJob.Name,
		JobImage:       &multiImage,
		JobCPU:         &multiCPU,
		JobDescription: &multiDescription,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, multiRequest)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with multiple parameters failed", err, multiRequest)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with multiple parameters: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	// Test 11: Patch job with run immediately
	log.Println("\nTest 11: Patch job with run immediately")
	newRunImmediately := true
	patchRequest11 := domain.PatchJobRequest{
		ProjectID:         cfg.ProjectID,
		JobName:           testJob.Name,
		JobRunImmediately: &newRunImmediately,
	}

	operation, err = jobs.PatchJob(cfg.ProjectID, testJob.Name, patchRequest11)
	if err != nil {
		logErrorWithRequestBody("Warning: Patch job with run immediately failed", err, patchRequest11)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully patched job with run immediately: %s, Operation ID: %s", testJob.Name, operation.ResourceID)
		} else {
			log.Printf("Warning: Operation validation failed for patch operation on job: %s", testJob.Name)
		}
	}

	log.Println("✓ Patch job test completed")
}
