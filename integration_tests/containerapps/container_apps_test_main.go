package main

import (
	"fmt"
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

func main() {
	cfg := config.LoadConfig()
	containerNameWithAutoDeploy := fmt.Sprintf("test-autodep-%s", cfg.ProjectID[:8])
	containerNameWithoutAutoDeploy := fmt.Sprintf("test-noautodep-%s", cfg.ProjectID[:8])

	// getListContainerApps(cfg)

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
	deleteContainerApp(cfg, containerNameWithAutoDeploy)

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
	deleteContainerApp(cfg, containerNameWithoutAutoDeploy)

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

	containerApp, err := ca.CreateContainerApp(request)
	log.Print(containerApp)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
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
