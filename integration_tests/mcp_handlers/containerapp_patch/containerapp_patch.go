package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/joho/godotenv"
)

// PatchTestCase represents a single patch test case with request and verification expectations
type PatchTestCase struct {
	Name          string
	ContainerName string
	Request       domain.PatchContainerAppRequest
	VerifyFunc    func(*domain.ContainerApp) error
}

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
	baseTimestamp := time.Now().Unix()

	// Cleanup all test-patch containerapps on exit
	defer func() {
		log.Println("=== Cleaning up all test-patch containerapps ===")
		cleanupAllTestPatchContainerApps(cfg)
	}()

	cleanupAllTestPatchContainerApps(cfg)

	// Define all patch test cases with unique container names
	patchTests := []PatchTestCase{
		{
			Name:          "Patch Basic Fields",
			ContainerName: fmt.Sprintf("test-patch-basic-%s-%d", cfg.ProjectID[:8], baseTimestamp),
			Request: domain.PatchContainerAppRequest{
				ProjectID:          cfg.ProjectID,
				ContainerAppPort:   intPtr(8081),
				ContainerAppImage:  strPtr("quickstart.cr.cloud.ru/restapi-go@sha256:d6bcdd96704c4db3ad176a975de5cfc403041422327fe54f9db89d5f249e0b87"),
				Description:        strPtr("Updated description for basic fields test"),
				PubliclyAccessible: boolPtr(false),
				Protocol:           strPtr("http_2"),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if len(app.Template.Containers) > 0 && app.Template.Containers[0].ContainerPort != 8081 {
					return fmt.Errorf("port verification failed: expected 8081, got %d", app.Template.Containers[0].ContainerPort)
				}
				if app.Description != "Updated description for basic fields test" {
					return fmt.Errorf("description verification failed: expected 'Updated description for basic fields test', got '%s'", app.Description)
				}
				if app.Template.Protocol != "http_2" {
					return fmt.Errorf("protocol verification failed: expected http_2, got %s", app.Template.Protocol)
				}
				return nil
			},
		},
		{
			Name:          "Patch Container Specific Fields",
			ContainerName: fmt.Sprintf("test-patch-container-%s-%d", cfg.ProjectID[:8], baseTimestamp+1),
			Request: domain.PatchContainerAppRequest{
				ProjectID:         cfg.ProjectID,
				ContainerAppImage: strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				ContainerAppPort:  intPtr(8080),
				CPU:               strPtr("0.2"),
				Timeout:           strPtr("30s"),
				IdleTimeout:       strPtr("300s"),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if len(app.Template.Containers) > 0 {
					if app.Template.Containers[0].Resources.CPU != "0.2" {
						return fmt.Errorf("CPU verification failed: expected 0.2, got %s", app.Template.Containers[0].Resources.CPU)
					}
				}
				if app.Template.Timeout != "30s" {
					return fmt.Errorf("timeout verification failed: expected 30s, got %s", app.Template.Timeout)
				}
				if app.Template.IdleTimeout != "300s" {
					return fmt.Errorf("idle timeout verification failed: expected 300s, got %s", app.Template.IdleTimeout)
				}
				return nil
			},
		},
		{
			Name:          "Patch Scaling Fields",
			ContainerName: fmt.Sprintf("test-patch-scaling-%s-%d", cfg.ProjectID[:8], baseTimestamp+2),
			Request: domain.PatchContainerAppRequest{
				ProjectID:         cfg.ProjectID,
				ContainerAppImage: strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				ContainerAppPort:  intPtr(8080),
				MinInstanceCount:  intPtr(1),
				MaxInstanceCount:  intPtr(5),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if app.Template.Scaling.MinInstanceCount != 1 {
					return fmt.Errorf("min instance count verification failed: expected 1, got %d", app.Template.Scaling.MinInstanceCount)
				}
				if app.Template.Scaling.MaxInstanceCount != 5 {
					return fmt.Errorf("max instance count verification failed: expected 5, got %d", app.Template.Scaling.MaxInstanceCount)
				}
				return nil
			},
		},
		{
			Name:          "Patch Auto-Deployments",
			ContainerName: fmt.Sprintf("test-patch-autodeploy-%s-%d", cfg.ProjectID[:8], baseTimestamp+3),
			Request: domain.PatchContainerAppRequest{
				ProjectID:              cfg.ProjectID,
				ContainerAppImage:      strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				ContainerAppPort:       intPtr(8080),
				AutoDeploymentsEnabled: boolPtr(true),
				AutoDeploymentsPattern: strPtr("latest"),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if app.Configuration.AutoDeployments.Enabled != true {
					return fmt.Errorf("auto-deployments enabled verification failed: expected true, got %v", app.Configuration.AutoDeployments.Enabled)
				}
				if app.Configuration.AutoDeployments.Pattern != "latest" {
					return fmt.Errorf("auto-deployments pattern verification failed: expected 'latest', got '%s'", app.Configuration.AutoDeployments.Pattern)
				}
				return nil
			},
		},
		{
			Name:          "Patch Environment Variables",
			ContainerName: fmt.Sprintf("test-patch-env-%s-%d", cfg.ProjectID[:8], baseTimestamp+4),
			Request: domain.PatchContainerAppRequest{
				ProjectID:            cfg.ProjectID,
				ContainerAppImage:    strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				ContainerAppPort:     intPtr(8080),
				EnvironmentVariables: strPtr("ENV_VAR1=value1;ENV_VAR2=value2;ENV_VAR3=value3"),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if len(app.Template.Containers) > 0 {
					envCount := len(app.Template.Containers[0].Env)
					if envCount != 3 {
						return fmt.Errorf("environment variables verification failed: expected 3 variables, got %d", envCount)
					}
					envMap := make(map[string]string)
					for _, env := range app.Template.Containers[0].Env {
						envMap[env.Name] = env.Value
					}
					if envMap["ENV_VAR1"] != "value1" {
						return fmt.Errorf("ENV_VAR1 verification failed: expected 'value1', got '%s'", envMap["ENV_VAR1"])
					}
					if envMap["ENV_VAR2"] != "value2" {
						return fmt.Errorf("ENV_VAR2 verification failed: expected 'value2', got '%s'", envMap["ENV_VAR2"])
					}
					if envMap["ENV_VAR3"] != "value3" {
						return fmt.Errorf("ENV_VAR3 verification failed: expected 'value3', got '%s'", envMap["ENV_VAR3"])
					}
				}
				return nil
			},
		},
		{
			Name:          "Patch Command and Args",
			ContainerName: fmt.Sprintf("test-patch-cmd-%s-%d", cfg.ProjectID[:8], baseTimestamp+5),
			Request: domain.PatchContainerAppRequest{
				ProjectID:         cfg.ProjectID,
				ContainerAppImage: strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				ContainerAppPort:  intPtr(8080),
				Command:           []string{"sh", "-c"},
				Args:              []string{"echo 'Hello from patched container!'"},
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if len(app.Template.Containers) > 0 {
					commandCount := len(app.Template.Containers[0].Command)
					if commandCount != 2 {
						return fmt.Errorf("command verification failed: expected 2 commands, got %d", commandCount)
					}
					argsCount := len(app.Template.Containers[0].Args)
					if argsCount != 1 {
						return fmt.Errorf("args verification failed: expected 1 arg, got %d", argsCount)
					}
				}
				return nil
			},
		},
		{
			Name:          "Patch Multiple Fields",
			ContainerName: fmt.Sprintf("test-patch-multi-%s-%d", cfg.ProjectID[:8], baseTimestamp+6),
			Request: domain.PatchContainerAppRequest{
				ProjectID:              cfg.ProjectID,
				ContainerAppPort:       intPtr(8082),
				ContainerAppImage:      strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				CPU:                    strPtr("0.3"),
				MinInstanceCount:       intPtr(1),
				MaxInstanceCount:       intPtr(4),
				Description:            strPtr("Multi-field patch test"),
				PubliclyAccessible:     boolPtr(true),
				Protocol:               strPtr("http_1"),
				Timeout:                strPtr("45s"),
				IdleTimeout:            strPtr("400s"),
				AutoDeploymentsEnabled: boolPtr(false),
				AutoDeploymentsPattern: strPtr("latest"),
				EnvironmentVariables:   strPtr("MULTI_VAR1=multi1;MULTI_VAR2=multi2"),
				Command:                []string{"node"},
				Args:                   []string{"server.js"},
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if len(app.Template.Containers) > 0 {
					if app.Template.Containers[0].ContainerPort != 8082 {
						return fmt.Errorf("port verification failed: expected 8082, got %d", app.Template.Containers[0].ContainerPort)
					}
					if app.Template.Containers[0].Resources.CPU != "0.3" {
						return fmt.Errorf("CPU verification failed: expected 0.3, got %s", app.Template.Containers[0].Resources.CPU)
					}
					if len(app.Template.Containers[0].Env) != 2 {
						return fmt.Errorf("environment variables verification failed: expected 2, got %d", len(app.Template.Containers[0].Env))
					}
					if len(app.Template.Containers[0].Command) != 1 {
						return fmt.Errorf("command verification failed: expected 1, got %d", len(app.Template.Containers[0].Command))
					}
					if len(app.Template.Containers[0].Args) != 1 {
						return fmt.Errorf("args verification failed: expected 1, got %d", len(app.Template.Containers[0].Args))
					}
				}
				if app.Description != "Multi-field patch test" {
					return fmt.Errorf("description verification failed: expected 'Multi-field patch test', got '%s'", app.Description)
				}
				if app.Configuration.Ingress.PubliclyAccessible != true {
					return fmt.Errorf("publicly accessible verification failed: expected true, got %v", app.Configuration.Ingress.PubliclyAccessible)
				}
				if app.Template.Protocol != "http_1" {
					return fmt.Errorf("protocol verification failed: expected http_1, got %s", app.Template.Protocol)
				}
				if app.Template.Timeout != "45s" {
					return fmt.Errorf("timeout verification failed: expected 45s, got %s", app.Template.Timeout)
				}
				if app.Template.IdleTimeout != "400s" {
					return fmt.Errorf("idle timeout verification failed: expected 400s, got %s", app.Template.IdleTimeout)
				}
				if app.Template.Scaling.MinInstanceCount != 1 {
					return fmt.Errorf("min instance count verification failed: expected 1, got %d", app.Template.Scaling.MinInstanceCount)
				}
				if app.Template.Scaling.MaxInstanceCount != 4 {
					return fmt.Errorf("max instance count verification failed: expected 4, got %d", app.Template.Scaling.MaxInstanceCount)
				}
				return nil
			},
		},
		{
			Name:          "Patch Partial Update",
			ContainerName: fmt.Sprintf("test-patch-partial-%s-%d", cfg.ProjectID[:8], baseTimestamp+7),
			Request: domain.PatchContainerAppRequest{
				ProjectID:         cfg.ProjectID,
				ContainerAppImage: strPtr("quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266"),
				Description:       strPtr("Partial update - only description changed"),
			},
			VerifyFunc: func(app *domain.ContainerApp) error {
				if app.Description != "Partial update - only description changed" {
					return fmt.Errorf("description verification failed: expected 'Partial update - only description changed', got '%s'", app.Description)
				}
				return nil
			},
		},
	}

	// Run integration tests
	log.Println("=== Running ContainerApp Patch Integration Tests in Parallel ===")

	// Create a WaitGroup to wait for all tests to complete
	var wg sync.WaitGroup
	wg.Add(len(patchTests))

	// Run each test in a separate goroutine
	for i, test := range patchTests {
		go func(testIndex int, testCase PatchTestCase) {
			defer wg.Done()
			log.Printf("\n--- Test %d/%d: %s (Container: %s) ---", testIndex+1, len(patchTests), testCase.Name, testCase.ContainerName)
			runPatchTest(cfg, testCase)
		}(i, test)
	}

	// Wait for all tests to complete
	wg.Wait()

	log.Println("=== All ContainerApp Patch Integration Tests Completed ===")
}

// runPatchTest executes a single patch test case
func runPatchTest(cfg *config.Config, test PatchTestCase) {
	// Create a test container app first
	log.Printf("Creating test container app: %s", test.ContainerName)
	createTestContainerApp(cfg, test.ContainerName)

	// Ensure cleanup
	defer func() {
		log.Printf("Cleaning up test container app: %s", test.ContainerName)
		deleteTestContainerApp(cfg, test.ContainerName)
	}()

	// Wait for container to be ready
	log.Printf("Waiting for container app to be ready: %s", test.ContainerName)
	waitForContainerAppReady(cfg, test.ContainerName, 30, 5*time.Second)

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Update request with container name
	test.Request.ContainerAppName = test.ContainerName

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, test.ContainerName, test.Request)
	if err != nil {
		log.Fatalf("Patch failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Wait for patch to complete
	log.Println("Waiting for patch to complete...")
	waitForContainerAppNotProcessing(cfg, test.ContainerName, 30, 5*time.Second)

	// Verify the patch
	// Add additional delay to ensure API is fully ready for next operation
	// Cloud.ru API needs significant time between patch operations to avoid 499 errors
	time.Sleep(15 * time.Second)
	containerApp := getContainerApp(cfg, test.ContainerName)

	// Run verification function
	if err := test.VerifyFunc(containerApp); err != nil {
		log.Fatalf("Verification failed: %v", err)
	}

	log.Printf("✓ %s test passed", test.Name)
}

// createTestContainerApp creates a test container app for patching
func createTestContainerApp(cfg *config.Config, name string) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	request := domain.CreateContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
		ContainerAppPort:       8080,
		ContainerAppImage:      "quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266",
		AutoDeploymentsEnabled: false,
		AutoDeploymentsPattern: "latest",
		Privileged:             false,
		IdleTimeout:            "600s",
		Timeout:                "60s",
		CPU:                    "0.1",
		MinInstanceCount:       0,
		MaxInstanceCount:       1,
		Description:            "Test container app for patch integration tests",
		PubliclyAccessible:     true,
		Protocol:               "http_1",
	}

	operation, err := ca.CreateContainerApp(request)
	if err != nil {
		log.Fatalf("Failed to create test container app: %v", err)
	}

	log.Printf("Test container app created: %+v", operation)
}

// deleteTestContainerApp deletes the test container app
func deleteTestContainerApp(cfg *config.Config, name string) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	_, err := ca.DeleteContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Printf("Warning: Failed to delete test container app: %v", err)
	} else {
		log.Printf("Test container app deleted successfully")
	}
}

// getContainerApp retrieves a container app for verification
func getContainerApp(cfg *config.Config, name string) *domain.ContainerApp {
	ca := cloudru.NewContainerAppsApplication(cfg)

	containerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Fatalf("Failed to get container app: %v", err)
	}

	return containerApp
}

// waitForContainerAppReady waits for the container app to be ready (not in "for_publish" status)
func waitForContainerAppReady(cfg *config.Config, name string, maxRetries int, retryInterval time.Duration) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	for i := 0; i < maxRetries; i++ {
		containerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to get container app: %v", i+1, maxRetries, err)
		} else {
			log.Printf("Attempt %d/%d: Container app status: %s", i+1, maxRetries, containerApp.Status)
			if containerApp.Status != "for_publish" {
				log.Printf("Container app is ready (status: %s)", containerApp.Status)
				return
			}
		}

		if i < maxRetries-1 {
			log.Printf("Waiting %v before next attempt...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Fatalf("Container app did not become ready after %d attempts", maxRetries)
}

// waitForContainerAppNotProcessing waits for the container app to not be in processing status
func waitForContainerAppNotProcessing(cfg *config.Config, name string, maxRetries int, retryInterval time.Duration) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	for i := 0; i < maxRetries; i++ {
		containerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to get container app: %v", i+1, maxRetries, err)
		} else {
			log.Printf("Attempt %d/%d: Container app status: %s", i+1, maxRetries, containerApp.Status)
			if containerApp.Status != "for_publish" && containerApp.Status != "for_patch" {
				log.Printf("Container app is ready (status: %s)", containerApp.Status)
				return
			}
		}

		if i < maxRetries-1 {
			log.Printf("Waiting %v before next attempt...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Fatalf("Container app did not become ready after %d attempts", maxRetries)
}

// Helper functions for creating pointer values
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// cleanupAllTestPatchContainerApps deletes all containerapps starting with "test-patch"
func cleanupAllTestPatchContainerApps(cfg *config.Config) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	// Get all containerapps
	containerApps, err := ca.GetListContainerApps(cfg.ProjectID)
	if err != nil {
		log.Printf("Warning: Failed to list containerapps for cleanup: %v", err)
		return
	}

	// Find and delete all containerapps starting with "test-patch"
	deletedCount := 0
	for _, app := range containerApps {
		if len(app.Name) >= 10 && app.Name[:10] == "test-patch" {
			log.Printf("Deleting test container app: %s", app.Name)
			_, err := ca.DeleteContainerApp(cfg.ProjectID, app.Name)
			if err != nil {
				log.Printf("Warning: Failed to delete test container app %s: %v", app.Name, err)
			} else {
				log.Printf("Successfully deleted test container app: %s", app.Name)
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		log.Printf("Cleanup completed: deleted %d test-patch containerapps", deletedCount)
	} else {
		log.Println("No test-patch containerapps found for cleanup")
	}
}
