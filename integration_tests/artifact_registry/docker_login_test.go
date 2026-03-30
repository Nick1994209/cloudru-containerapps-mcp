package main

import (
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	// Test docker login with a specific registry name
	testDockerLogin(cfg, "test-registry-1774870656")
}

func testDockerLogin(cfg *config.Config, registryName string) {
	dockerApp := application.NewDockerApplication(cfg)

	log.Printf("Testing Docker login with registry: %s...", registryName)
	registry, err := dockerApp.Login(registryName)
	if err != nil {
		log.Printf("Docker login error: %v", err)
	} else {
		log.Printf("Docker login success: logged into %s", registry)
	}
}
