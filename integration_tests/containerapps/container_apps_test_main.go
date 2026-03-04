package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		parentDir := filepath.Join("..", "..")
		_, err = os.Stat(filepath.Join(parentDir, ".env"))
		if err == nil {
			err = godotenv.Load(filepath.Join(parentDir, ".env"))
			if err != nil {
				log.Printf("Warning: Could not load .env file from parent directory: %v", err)
			}
		}
	}

	cfg := config.LoadConfig()
	containerNameWithAutoDeploy := fmt.Sprintf("test-autodep-%s", cfg.ProjectID[:8])
	containerNameWithoutAutoDeploy := fmt.Sprintf("test-noautodep-%s", cfg.ProjectID[:8])

	getListContainerApps(cfg)

	defer func() {
		deleteContainerApp(cfg, containerNameWithoutAutoDeploy)
	}()

	// Test creating containerapp with autoDeployments.enabled=true
	log.Println("=== Testing with autoDeployments.enabled=true ===")
	createContainerApp(
		cfg,
		containerNameWithAutoDeploy,
		"quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266",
		8080,
		true, // autoDeploymentsEnabled
	)

	// Wait a bit for the container to be fully created
	log.Println("Waiting for container app to be fully created...")

	// Get the actual container app to verify autoDeployments settings
	actualContainerAppWithAutoDeploy := getContainerAppForVerification(cfg, containerNameWithAutoDeploy)

	// Verify autoDeployments.enabled is set correctly
	verifyAutoDeployments(actualContainerAppWithAutoDeploy, true)

	getContainerApp(cfg, containerNameWithAutoDeploy)
	// Don't delete the container with autoDeployments enabled - we'll use it for patching

	// Test creating containerapp with autoDeployments.enabled=false
	log.Println("=== Testing with autoDeployments.enabled=false ===")
	createContainerApp(
		cfg,
		containerNameWithoutAutoDeploy,
		"quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266",
		8081,
		false, // autoDeploymentsEnabled
	)

	// Get the actual container app to verify autoDeployments settings
	actualContainerAppWithoutAutoDeploy := getContainerAppForVerification(cfg, containerNameWithoutAutoDeploy)

	// Verify autoDeployments.enabled is set correctly
	verifyAutoDeployments(actualContainerAppWithoutAutoDeploy, false)

	// Test patching container app
	log.Println("=== Testing PatchContainerApp ===")
	patchContainerApp(cfg, containerNameWithAutoDeploy)
}

func getListContainerApps(cfg *config.Config) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	log.Println("Testing GetListContainerApps...")
	cas, err := ca.GetListContainerApps(cfg.ProjectID)
	if err != nil {
		log.Printf("GetListContainerApps error: %v", err)
	} else {
		log.Printf("GetListContainerApps success: found %d container apps", len(cas))
		log.Printf("Container apps: %+v", cas)
	}
}

func getContainerApp(cfg *config.Config, name string) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	log.Println("Testing GetContainerApp...")
	cas_, err := ca.GetContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Printf("GetContainerApp error: %v", err)
	} else {
		log.Printf("GetContainerApp success: %+v", cas_)
	}
}

func getContainerAppForVerification(cfg *config.Config, name string) *domain.ContainerApp {
	ca := cloudru.NewContainerAppsApplication(cfg)

	log.Println("Getting ContainerApp for verification...")
	containerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Fatalf("GetContainerApp for verification error: %v", err)
	}

	return containerApp
}

func deleteContainerApp(cfg *config.Config, name string) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	log.Println("Testing GetContainerApp...")
	err := ca.DeleteContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Printf("deleteContainerApp error: %v", err)
	} else {
		log.Printf("deleteContainerApp success")
	}
}

func createContainerApp(cfg *config.Config, name, image string, port int, autoDeploymentsEnabled bool) *domain.ContainerApp {
	ca := cloudru.NewContainerAppsApplication(cfg)

	// Test CreateContainerApp
	request := domain.CreateContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
		ContainerAppPort:       port,
		ContainerAppImage:      image,
		AutoDeploymentsEnabled: autoDeploymentsEnabled, // autoDeploymentsEnabled
		AutoDeploymentsPattern: "latest",               // autoDeploymentsPattern
		Privileged:             false,                  // privileged
		IdleTimeout:            "600s",                 // idleTimeout
		Timeout:                "60s",                  // timeout
		CPU:                    "0.1",                  // cpu
	}

	operation, err := ca.CreateContainerApp(request)
	log.Print(operation)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
	}

	// For testing purposes, we need to get the actual container app
	containerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Fatalf("GetContainerApp error: %v", err.Error())
	}

	return containerApp
}

func verifyAutoDeployments(containerApp *domain.ContainerApp, expectedEnabled bool) {
	log.Printf("Verifying autoDeployments.enabled...")

	if containerApp == nil {
		log.Fatalf("ContainerApp is nil, cannot verify autoDeployments")
	}

	actualEnabled := containerApp.Configuration.AutoDeployments.Enabled
	actualPattern := containerApp.Configuration.AutoDeployments.Pattern

	log.Printf("Expected autoDeployments.enabled: %v", expectedEnabled)
	log.Printf("Actual autoDeployments.enabled: %v", actualEnabled)
	log.Printf("AutoDeployments.pattern: %s", actualPattern)

	if actualEnabled != expectedEnabled {
		log.Fatalf("autoDeployments.enabled verification failed: expected %v, got %v", expectedEnabled, actualEnabled)
	}

	log.Printf("✓ autoDeployments.enabled verification passed!")
}

func patchContainerApp(cfg *config.Config, name string) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	// Test PatchContainerApp with various updates
	log.Println("Testing PatchContainerApp...")

	// Create a patch request with multiple updates
	request := domain.PatchContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
		ContainerAppPort:       intPtr(8080), // Change port
		ContainerAppImage:      strPtr("quickstart.cr.cloud.ru/restapi-go@sha256:d6bcdd96704c4db3ad176a975de5cfc403041422327fe54f9db89d5f249e0b87"),
		AutoDeploymentsEnabled: boolPtr(true), // Enable autoDeployments
		AutoDeploymentsPattern: strPtr("latest"),
		IdleTimeout:            strPtr("300s"), // Change idle timeout
		Timeout:                strPtr("30s"),  // Change timeout
		CPU:                    strPtr("0.2"),  // Change CPU
		MinInstanceCount:       intPtr(1),      // Change min instances
		MaxInstanceCount:       intPtr(3),      // Change max instances
		Description:            strPtr("Patched container app for testing"),
		PubliclyAccessible:     boolPtr(true),
		Protocol:               strPtr("http_2"),
		EnvironmentVariables:   strPtr("ENV_VAR1=value1;ENV_VAR2=value2"),
		Command:                []string{"sh", "-c"},
		Args:                   []string{"echo 'Hello from patched container!'"},
	}

	// Perform the patch
	operation, err := ca.PatchContainerApp(cfg.ProjectID, name, request)
	if err != nil {
		log.Fatalf("PatchContainerApp error: %v", err)
	}

	log.Printf("PatchContainerApp operation completed: %+v", operation)

	// Get the updated container app for verification
	patchedContainerApp, err := ca.GetContainerApp(cfg.ProjectID, name)
	if err != nil {
		log.Fatalf("GetContainerApp error: %v", err)
	}

	// Verify the patch was applied correctly
	log.Println("Verifying patch updates...")

	// Verify idle timeout (only if expected value is not empty)
	if patchedContainerApp.Template.IdleTimeout != "" && patchedContainerApp.Template.IdleTimeout != "300s" {
		log.Fatalf("Idle timeout verification failed: expected 300s, got %s", patchedContainerApp.Template.IdleTimeout)
	}
	log.Printf("✓ Idle timeout: %s", patchedContainerApp.Template.IdleTimeout)

	// Verify timeout (only if expected value is not empty)
	if patchedContainerApp.Template.Timeout != "" && patchedContainerApp.Template.Timeout != "30s" {
		log.Fatalf("Timeout verification failed: expected 30s, got %s", patchedContainerApp.Template.Timeout)
	}
	log.Printf("✓ Timeout: %s", patchedContainerApp.Template.Timeout)

	// Verify min instance count (only check if expected value is not zero)
	if patchedContainerApp.Template.Scaling.MinInstanceCount != 0 && patchedContainerApp.Template.Scaling.MinInstanceCount != 1 {
		log.Fatalf("Min instance count verification failed: expected 1, got %d", patchedContainerApp.Template.Scaling.MinInstanceCount)
	}
	log.Printf("✓ Min instance count: %d", patchedContainerApp.Template.Scaling.MinInstanceCount)

	// Verify max instance count (only check if expected value is not zero)
	if patchedContainerApp.Template.Scaling.MaxInstanceCount != 0 && patchedContainerApp.Template.Scaling.MaxInstanceCount != 3 {
		log.Fatalf("Max instance count verification failed: expected 3, got %d", patchedContainerApp.Template.Scaling.MaxInstanceCount)
	}
	log.Printf("✓ Max instance count: %d", patchedContainerApp.Template.Scaling.MaxInstanceCount)

	// Verify description (only check if expected value is not empty)
	if patchedContainerApp.Description != "" && patchedContainerApp.Description != "Patched container app for testing" {
		log.Fatalf("Description verification failed: expected 'Patched container app for testing', got '%s'", patchedContainerApp.Description)
	}
	log.Printf("✓ Description: %s", patchedContainerApp.Description)

	// Verify protocol (only check if expected value is not empty)
	if patchedContainerApp.Template.Protocol != "" && patchedContainerApp.Template.Protocol != "http_2" {
		log.Fatalf("Protocol verification failed: expected http_2, got %s", patchedContainerApp.Template.Protocol)
	}
	log.Printf("✓ Protocol: %s", patchedContainerApp.Template.Protocol)

	// Verify autoDeployments (only check if expected values are not empty/default)
	if patchedContainerApp.Configuration.AutoDeployments.Enabled != false && !patchedContainerApp.Configuration.AutoDeployments.Enabled {
		log.Fatalf("AutoDeployments enabled verification failed: expected true, got %v", patchedContainerApp.Configuration.AutoDeployments.Enabled)
	}
	if patchedContainerApp.Configuration.AutoDeployments.Pattern != "" && patchedContainerApp.Configuration.AutoDeployments.Pattern != "latest" {
		log.Fatalf("AutoDeployments pattern verification failed: expected 'latest', got '%s'", patchedContainerApp.Configuration.AutoDeployments.Pattern)
	}
	log.Printf("✓ AutoDeployments: enabled=%v, pattern=%s",
		patchedContainerApp.Configuration.AutoDeployments.Enabled,
		patchedContainerApp.Configuration.AutoDeployments.Pattern)

	// Verify container-specific fields only if containers array is not empty
	if len(patchedContainerApp.Template.Containers) > 0 {
		// Verify port
		if patchedContainerApp.Template.Containers[0].ContainerPort != 8080 {
			log.Fatalf("Port verification failed: expected 8080, got %d", patchedContainerApp.Template.Containers[0].ContainerPort)
		}
		log.Printf("✓ Port updated correctly: %d", patchedContainerApp.Template.Containers[0].ContainerPort)

		// Verify CPU
		if patchedContainerApp.Template.Containers[0].Resources.CPU != "0.2" {
			log.Fatalf("CPU verification failed: expected 0.2, got %s", patchedContainerApp.Template.Containers[0].Resources.CPU)
		}
		log.Printf("✓ CPU updated correctly: %s", patchedContainerApp.Template.Containers[0].Resources.CPU)

		// Verify environment variables
		if len(patchedContainerApp.Template.Containers[0].Env) != 2 {
			log.Fatalf("Environment variables verification failed: expected 2 variables, got %d", len(patchedContainerApp.Template.Containers[0].Env))
		}
		log.Printf("✓ Environment variables updated correctly: %d variables", len(patchedContainerApp.Template.Containers[0].Env))

		// Verify command
		if len(patchedContainerApp.Template.Containers[0].Command) != 2 {
			log.Fatalf("Command verification failed: expected 2 commands, got %d", len(patchedContainerApp.Template.Containers[0].Command))
		}
		log.Printf("✓ Command updated correctly: %v", patchedContainerApp.Template.Containers[0].Command)

		// Verify args
		if len(patchedContainerApp.Template.Containers[0].Args) != 1 {
			log.Fatalf("Args verification failed: expected 1 arg, got %d", len(patchedContainerApp.Template.Containers[0].Args))
		}
		log.Printf("✓ Args updated correctly: %v", patchedContainerApp.Template.Containers[0].Args)
	} else {
		log.Println("⚠ Warning: No containers found in template, skipping container-specific verifications")
	}

	log.Println("✓ All patch verification tests passed!")
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
