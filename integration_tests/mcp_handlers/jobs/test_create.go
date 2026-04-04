package main

import (
	"log"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// testCreateJob tests creating a job with various parameters
func testCreateJob(cfg *config.Config) []TestJob {
	log.Println("\n--- Test: Create Job ---")

	jobs := cloudru.NewJobsApplication(cfg)
	var createdJobs []TestJob

	// Add a small delay to ensure any previous operations are complete
	time.Sleep(2 * time.Second)

	// Test 1: Create job with minimal parameters
	log.Println("\nTest 1: Create job with minimal parameters")
	minimalJobName := "test-job-min-" + time.Now().Format("150405")
	minimalRequest := domain.CreateJobRequest{
		ProjectID:           cfg.ProjectID,
		JobName:             minimalJobName,
		JobImage:            "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2",
		JobPrivileged:       false,
		JobCPU:              "0.1",
		JobRetryCount:       1,
		JobExecutionTimeout: 60,
	}

	operation, err := jobs.CreateJob(minimalRequest)
	if err != nil {
		logErrorWithRequestBody("Warning: Create job with minimal parameters failed", err, minimalRequest)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully created job with minimal parameters: %s, Operation ID: %s", minimalJobName, operation.ResourceID)
			createdJobs = append(createdJobs, TestJob{Name: minimalJobName, ID: operation.ResourceID})
		} else {
			log.Printf("Warning: Operation validation failed for job: %s", minimalJobName)
		}
	}

	// Test 2: Create job with all parameters
	log.Println("\nTest 2: Create job with all parameters")
	fullJobName := "test-job-full-" + time.Now().Format("150405")
	fullRequest := domain.CreateJobRequest{
		ProjectID:               cfg.ProjectID,
		JobName:                 fullJobName,
		JobImage:                "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2",
		JobPrivileged:           true,
		JobCPU:                  "0.2",
		JobDescription:          "Test job with all parameters",
		JobEnvironmentVariables: "ENV1='value1';ENV2='value2'",
		JobCommand:              []string{"/bin/sh", "-c"},
		JobArgs:                 []string{"echo", "Hello World"},
		JobRetryCount:           3,
		JobExecutionTimeout:     600,
		JobRunImmediately:       true,
	}

	operation, err = jobs.CreateJob(fullRequest)
	if err != nil {
		logErrorWithRequestBody("Warning: Create job with all parameters failed", err, fullRequest)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully created job with all parameters: %s, Operation ID: %s", fullJobName, operation.ResourceID)
			createdJobs = append(createdJobs, TestJob{Name: fullJobName, ID: operation.ResourceID})
		} else {
			log.Printf("Warning: Operation validation failed for job: %s", fullJobName)
		}
	}

	// Test 3: Create job with command and args only
	log.Println("\nTest 3: Create job with command and args only")
	cmdJobName := "test-job-cmd-" + time.Now().Format("150405")
	cmdRequest := domain.CreateJobRequest{
		ProjectID:           cfg.ProjectID,
		JobName:             cmdJobName,
		JobImage:            "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2",
		JobPrivileged:       false,
		JobCPU:              "0.1",
		JobCommand:          []string{"/bin/sh"},
		JobArgs:             []string{"-c", "echo 'Test command'"},
		JobRetryCount:       1,
		JobExecutionTimeout: 60,
	}

	operation, err = jobs.CreateJob(cmdRequest)
	if err != nil {
		logErrorWithRequestBody("Warning: Create job with command and args failed", err, cmdRequest)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully created job with command and args: %s, Operation ID: %s", cmdJobName, operation.ResourceID)
			createdJobs = append(createdJobs, TestJob{Name: cmdJobName, ID: operation.ResourceID})
		} else {
			log.Printf("Warning: Operation validation failed for job: %s", cmdJobName)
		}
	}

	// Test 4: Create job with environment variables only
	log.Println("\nTest 4: Create job with environment variables only")
	envJobName := "test-job-env-" + time.Now().Format("150405")
	envRequest := domain.CreateJobRequest{
		ProjectID:               cfg.ProjectID,
		JobName:                 envJobName,
		JobImage:                "nvkorolkov-public.cr.cloud.ru/job-did-you-know:v0.0.2",
		JobPrivileged:           false,
		JobCPU:                  "0.1",
		JobEnvironmentVariables: "API_KEY='test123';DB_HOST='localhost';DB_PORT='5432'",
		JobRetryCount:           1,
		JobExecutionTimeout:     60,
	}

	operation, err = jobs.CreateJob(envRequest)
	if err != nil {
		logErrorWithRequestBody("Warning: Create job with environment variables failed", err, envRequest)
	} else {
		if validateOperation(operation) {
			log.Printf("Successfully created job with environment variables: %s, Operation ID: %s", envJobName, operation.ResourceID)
			createdJobs = append(createdJobs, TestJob{Name: envJobName, ID: operation.ResourceID})
		} else {
			log.Printf("Warning: Operation validation failed for job: %s", envJobName)
		}
	}

	log.Printf("✓ Create job test completed, created %d job(s)", len(createdJobs))
	return createdJobs
}
