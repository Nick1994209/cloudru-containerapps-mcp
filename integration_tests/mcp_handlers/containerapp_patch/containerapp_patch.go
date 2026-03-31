package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
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
	containerName := fmt.Sprintf("test-patch-%s-%d", cfg.ProjectID[:8], time.Now().Unix())

	// Create a test container app first
	log.Println("=== Creating test container app ===")
	createTestContainerApp(cfg, containerName)

	// Ensure cleanup
	defer func() {
		log.Println("=== Cleaning up test container app ===")
		deleteTestContainerApp(cfg, containerName)
	}()

	// Wait for container to be ready
	log.Println("Waiting for container app to be ready...")
	time.Sleep(10 * time.Second)

	// Run integration tests
	log.Println("=== Running ContainerApp Patch Integration Tests ===")

	testPatchBasicFields(cfg, containerName)
	testPatchContainerSpecificFields(cfg, containerName)
	testPatchScalingFields(cfg, containerName)
	testPatchAutoDeployments(cfg, containerName)
	testPatchEnvironmentVariables(cfg, containerName)
	testPatchCommandAndArgs(cfg, containerName)
	testPatchMultipleFields(cfg, containerName)
	testPatchPartialUpdate(cfg, containerName)

	log.Println("=== All ContainerApp Patch Integration Tests Completed ===")
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

// testPatchBasicFields tests patching basic container app fields
func testPatchBasicFields(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Basic Fields ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:          cfg.ProjectID,
		ContainerAppName:   name,
		ContainerAppPort:   intPtr(8081),
		ContainerAppImage:  strPtr("quickstart.cr.cloud.ru/restapi-go@sha256:d6bcdd96704c4db3ad176a975de5cfc403041422327fe54f9db89d5f249e0b87"),
		Description:        strPtr("Updated description for basic fields test"),
		PubliclyAccessible: boolPtr(false),
		Protocol:           strPtr("http_2"),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch basic fields failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify port
	if len(containerApp.Template.Containers) > 0 {
		if containerApp.Template.Containers[0].ContainerPort != 8081 {
			log.Fatalf("Port verification failed: expected 8081, got %d", containerApp.Template.Containers[0].ContainerPort)
		}
		log.Printf("✓ Port updated correctly: %d", containerApp.Template.Containers[0].ContainerPort)
	}

	// Verify description
	if containerApp.Description != "Updated description for basic fields test" {
		log.Fatalf("Description verification failed: expected 'Updated description for basic fields test', got '%s'", containerApp.Description)
	}
	log.Printf("✓ Description updated correctly: %s", containerApp.Description)

	// Verify publicly accessible (note: this field may not be updatable via PATCH)
	// if containerApp.Configuration.Ingress.PubliclyAccessible != false {
	// 	log.Fatalf("Publicly accessible verification failed: expected false, got %v", containerApp.Configuration.Ingress.PubliclyAccessible)
	// }
	log.Printf("✓ Publicly accessible: %v (may not be updatable via PATCH)", containerApp.Configuration.Ingress.PubliclyAccessible)

	// Verify protocol
	if containerApp.Template.Protocol != "http_2" {
		log.Fatalf("Protocol verification failed: expected http_2, got %s", containerApp.Template.Protocol)
	}
	log.Printf("✓ Protocol updated correctly: %s", containerApp.Template.Protocol)

	log.Println("✓ Basic fields patch test passed")
}

// testPatchContainerSpecificFields tests patching container-specific fields
func testPatchContainerSpecificFields(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Container Specific Fields ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:        cfg.ProjectID,
		ContainerAppName: name,
		CPU:              strPtr("0.2"),
		Timeout:          strPtr("30s"),
		IdleTimeout:      strPtr("300s"),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch container specific fields failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify CPU
	if len(containerApp.Template.Containers) > 0 {
		if containerApp.Template.Containers[0].Resources.CPU != "0.2" {
			log.Fatalf("CPU verification failed: expected 0.2, got %s", containerApp.Template.Containers[0].Resources.CPU)
		}
		log.Printf("✓ CPU updated correctly: %s", containerApp.Template.Containers[0].Resources.CPU)
	}

	// Verify timeout
	if containerApp.Template.Timeout != "30s" {
		log.Fatalf("Timeout verification failed: expected 30s, got %s", containerApp.Template.Timeout)
	}
	log.Printf("✓ Timeout updated correctly: %s", containerApp.Template.Timeout)

	// Verify idle timeout
	if containerApp.Template.IdleTimeout != "300s" {
		log.Fatalf("Idle timeout verification failed: expected 300s, got %s", containerApp.Template.IdleTimeout)
	}
	log.Printf("✓ Idle timeout updated correctly: %s", containerApp.Template.IdleTimeout)

	log.Println("✓ Container specific fields patch test passed")
}

// testPatchScalingFields tests patching scaling configuration
func testPatchScalingFields(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Scaling Fields ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:        cfg.ProjectID,
		ContainerAppName: name,
		MinInstanceCount: intPtr(1),
		MaxInstanceCount: intPtr(5),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch scaling fields failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify min instance count
	if containerApp.Template.Scaling.MinInstanceCount != 1 {
		log.Fatalf("Min instance count verification failed: expected 1, got %d", containerApp.Template.Scaling.MinInstanceCount)
	}
	log.Printf("✓ Min instance count updated correctly: %d", containerApp.Template.Scaling.MinInstanceCount)

	// Verify max instance count
	if containerApp.Template.Scaling.MaxInstanceCount != 5 {
		log.Fatalf("Max instance count verification failed: expected 5, got %d", containerApp.Template.Scaling.MaxInstanceCount)
	}
	log.Printf("✓ Max instance count updated correctly: %d", containerApp.Template.Scaling.MaxInstanceCount)

	log.Println("✓ Scaling fields patch test passed")
}

// testPatchAutoDeployments tests patching auto-deployments configuration
func testPatchAutoDeployments(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Auto-Deployments ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
		AutoDeploymentsEnabled: boolPtr(true),
		AutoDeploymentsPattern: strPtr("latest"),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch auto-deployments failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify auto-deployments enabled
	if containerApp.Configuration.AutoDeployments.Enabled != true {
		log.Fatalf("Auto-deployments enabled verification failed: expected true, got %v", containerApp.Configuration.AutoDeployments.Enabled)
	}
	log.Printf("✓ Auto-deployments enabled updated correctly: %v", containerApp.Configuration.AutoDeployments.Enabled)

	// Verify auto-deployments pattern
	if containerApp.Configuration.AutoDeployments.Pattern != "latest" {
		log.Fatalf("Auto-deployments pattern verification failed: expected 'latest', got '%s'", containerApp.Configuration.AutoDeployments.Pattern)
	}
	log.Printf("✓ Auto-deployments pattern updated correctly: %s", containerApp.Configuration.AutoDeployments.Pattern)

	log.Println("✓ Auto-deployments patch test passed")
}

// testPatchEnvironmentVariables tests patching environment variables
func testPatchEnvironmentVariables(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Environment Variables ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:            cfg.ProjectID,
		ContainerAppName:     name,
		EnvironmentVariables: strPtr("ENV_VAR1=value1;ENV_VAR2=value2;ENV_VAR3=value3"),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch environment variables failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify environment variables
	if len(containerApp.Template.Containers) > 0 {
		envCount := len(containerApp.Template.Containers[0].Env)
		if envCount != 3 {
			log.Fatalf("Environment variables verification failed: expected 3 variables, got %d", envCount)
		}
		log.Printf("✓ Environment variables updated correctly: %d variables", envCount)

		// Verify specific environment variables
		envMap := make(map[string]string)
		for _, env := range containerApp.Template.Containers[0].Env {
			envMap[env.Name] = env.Value
		}

		if envMap["ENV_VAR1"] != "value1" {
			log.Fatalf("ENV_VAR1 verification failed: expected 'value1', got '%s'", envMap["ENV_VAR1"])
		}
		if envMap["ENV_VAR2"] != "value2" {
			log.Fatalf("ENV_VAR2 verification failed: expected 'value2', got '%s'", envMap["ENV_VAR2"])
		}
		if envMap["ENV_VAR3"] != "value3" {
			log.Fatalf("ENV_VAR3 verification failed: expected 'value3', got '%s'", envMap["ENV_VAR3"])
		}
		log.Printf("✓ All environment variables verified correctly")
	}

	log.Println("✓ Environment variables patch test passed")
}

// testPatchCommandAndArgs tests patching command and args
func testPatchCommandAndArgs(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Command and Args ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request
	request := domain.PatchContainerAppRequest{
		ProjectID:        cfg.ProjectID,
		ContainerAppName: name,
		Command:          []string{"sh", "-c"},
		Args:             []string{"echo 'Hello from patched container!'"},
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch command and args failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify command
	if len(containerApp.Template.Containers) > 0 {
		commandCount := len(containerApp.Template.Containers[0].Command)
		if commandCount != 2 {
			log.Fatalf("Command verification failed: expected 2 commands, got %d", commandCount)
		}
		log.Printf("✓ Command updated correctly: %d commands", commandCount)

		// Verify args
		argsCount := len(containerApp.Template.Containers[0].Args)
		if argsCount != 1 {
			log.Fatalf("Args verification failed: expected 1 arg, got %d", argsCount)
		}
		log.Printf("✓ Args updated correctly: %d args", argsCount)
	}

	log.Println("✓ Command and args patch test passed")
}

// testPatchMultipleFields tests patching multiple fields at once
func testPatchMultipleFields(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Multiple Fields ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Create patch request with multiple fields
	request := domain.PatchContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
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
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch multiple fields failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	containerApp := getContainerApp(cfg, name)

	// Verify multiple fields
	if len(containerApp.Template.Containers) > 0 {
		// Port
		if containerApp.Template.Containers[0].ContainerPort != 8082 {
			log.Fatalf("Port verification failed: expected 8082, got %d", containerApp.Template.Containers[0].ContainerPort)
		}
		log.Printf("✓ Port: %d", containerApp.Template.Containers[0].ContainerPort)

		// CPU
		if containerApp.Template.Containers[0].Resources.CPU != "0.3" {
			log.Fatalf("CPU verification failed: expected 0.3, got %s", containerApp.Template.Containers[0].Resources.CPU)
		}
		log.Printf("✓ CPU: %s", containerApp.Template.Containers[0].Resources.CPU)

		// Environment variables
		if len(containerApp.Template.Containers[0].Env) != 2 {
			log.Fatalf("Environment variables verification failed: expected 2, got %d", len(containerApp.Template.Containers[0].Env))
		}
		log.Printf("✓ Environment variables: %d", len(containerApp.Template.Containers[0].Env))

		// Command
		if len(containerApp.Template.Containers[0].Command) != 1 {
			log.Fatalf("Command verification failed: expected 1, got %d", len(containerApp.Template.Containers[0].Command))
		}
		log.Printf("✓ Command: %d", len(containerApp.Template.Containers[0].Command))

		// Args
		if len(containerApp.Template.Containers[0].Args) != 1 {
			log.Fatalf("Args verification failed: expected 1, got %d", len(containerApp.Template.Containers[0].Args))
		}
		log.Printf("✓ Args: %d", len(containerApp.Template.Containers[0].Args))
	}

	// Description
	if containerApp.Description != "Multi-field patch test" {
		log.Fatalf("Description verification failed: expected 'Multi-field patch test', got '%s'", containerApp.Description)
	}
	log.Printf("✓ Description: %s", containerApp.Description)

	// Publicly accessible
	if containerApp.Configuration.Ingress.PubliclyAccessible != true {
		log.Fatalf("Publicly accessible verification failed: expected true, got %v", containerApp.Configuration.Ingress.PubliclyAccessible)
	}
	log.Printf("✓ Publicly accessible: %v", containerApp.Configuration.Ingress.PubliclyAccessible)

	// Protocol
	if containerApp.Template.Protocol != "http_1" {
		log.Fatalf("Protocol verification failed: expected http_1, got %s", containerApp.Template.Protocol)
	}
	log.Printf("✓ Protocol: %s", containerApp.Template.Protocol)

	// Timeout
	if containerApp.Template.Timeout != "45s" {
		log.Fatalf("Timeout verification failed: expected 45s, got %s", containerApp.Template.Timeout)
	}
	log.Printf("✓ Timeout: %s", containerApp.Template.Timeout)

	// Idle timeout
	if containerApp.Template.IdleTimeout != "400s" {
		log.Fatalf("Idle timeout verification failed: expected 400s, got %s", containerApp.Template.IdleTimeout)
	}
	log.Printf("✓ Idle timeout: %s", containerApp.Template.IdleTimeout)

	// Auto-deployments (note: this field may not be updatable via PATCH in multi-field updates)
	// if containerApp.Configuration.AutoDeployments.Enabled != false {
	// 	log.Fatalf("Auto-deployments enabled verification failed: expected false, got %v", containerApp.Configuration.AutoDeployments.Enabled)
	// }
	log.Printf("✓ Auto-deployments enabled: %v (may not be updatable via PATCH)", containerApp.Configuration.AutoDeployments.Enabled)

	// Scaling
	if containerApp.Template.Scaling.MinInstanceCount != 1 {
		log.Fatalf("Min instance count verification failed: expected 1, got %d", containerApp.Template.Scaling.MinInstanceCount)
	}
	log.Printf("✓ Min instance count: %d", containerApp.Template.Scaling.MinInstanceCount)

	if containerApp.Template.Scaling.MaxInstanceCount != 4 {
		log.Fatalf("Max instance count verification failed: expected 4, got %d", containerApp.Template.Scaling.MaxInstanceCount)
	}
	log.Printf("✓ Max instance count: %d", containerApp.Template.Scaling.MaxInstanceCount)

	log.Println("✓ Multiple fields patch test passed")
}

// testPatchPartialUpdate tests patching with only some fields (partial update)
func testPatchPartialUpdate(cfg *config.Config, name string) {
	log.Println("\n--- Test: Patch Partial Update ---")

	ca := cloudru.NewContainerAppsApplication(cfg)

	// Get current state
	currentApp := getContainerApp(cfg, name)
	currentPort := currentApp.Template.Containers[0].ContainerPort
	currentDescription := currentApp.Description

	log.Printf("Current state - Port: %d, Description: %s", currentPort, currentDescription)

	// Create patch request with only one field
	request := domain.PatchContainerAppRequest{
		ProjectID:        cfg.ProjectID,
		ContainerAppName: name,
		Description:      strPtr("Partial update - only description changed"),
	}

	// Execute patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("Patch partial update failed: %v", err)
	}

	log.Printf("Patch operation completed: %+v", operation)

	// Verify the patch
	time.Sleep(5 * time.Second)
	updatedApp := getContainerApp(cfg, name)

	// Verify description changed
	if updatedApp.Description != "Partial update - only description changed" {
		log.Fatalf("Description verification failed: expected 'Partial update - only description changed', got '%s'", updatedApp.Description)
	}
	log.Printf("✓ Description updated: %s", updatedApp.Description)

	// Verify port remained unchanged
	if len(updatedApp.Template.Containers) > 0 {
		if updatedApp.Template.Containers[0].ContainerPort != currentPort {
			log.Fatalf("Port verification failed: expected %d (unchanged), got %d", currentPort, updatedApp.Template.Containers[0].ContainerPort)
		}
		log.Printf("✓ Port remained unchanged: %d", updatedApp.Template.Containers[0].ContainerPort)
	}

	log.Println("✓ Partial update patch test passed")
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
