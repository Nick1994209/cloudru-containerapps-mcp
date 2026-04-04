package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from the project root directory
	_, err := os.Stat(".env")
	if err == nil {
		// .env exists in current directory
		err = godotenv.Load()
		if err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	} else {
		// Try to load from parent directory (project root)
		parentDir := filepath.Join("..", "..", "..")
		_, err = os.Stat(filepath.Join(parentDir, ".env"))
		if err == nil {
			err = godotenv.Load(filepath.Join(parentDir, ".env"))
			if err != nil {
				log.Printf("Warning: Could not load .env file from parent directory: %v", err)
			}
		}
	}

	cfg := config.LoadConfig()

	// Run integration tests
	log.Println("=== Running Jobs Integration Tests ===")

	// Cleanup: Delete all test jobs (defer ensures this always runs)
	defer cleanupTestJobs(cfg)

	// // Create jobs and collect them
	createdJobs := testCreateJob(cfg)

	testListJobs(cfg, createdJobs)

	// Patch the first created job
	if len(createdJobs) > 0 {
		testGetJob(cfg, createdJobs[0])
		testExecuteJob(cfg, createdJobs[0])
		testListJobExecutions(cfg, createdJobs[0])
		testPatchJob(cfg, createdJobs[0])
		testDeleteJob(cfg, createdJobs)
	}

	log.Println("=== All Jobs Integration Tests Completed ===")
}
