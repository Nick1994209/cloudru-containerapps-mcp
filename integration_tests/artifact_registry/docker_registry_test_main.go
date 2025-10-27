package main

import (
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	getListDockerRegistries(cfg)
	createDockerRegistry(cfg, "test-registry", false)
}

func getListDockerRegistries(cfg *config.Config) {
	dr := cloudru.NewArtifactRegistryApplication(cfg)

	log.Println("Testing GetListDockerRegistries...")
	registries, err := dr.GetListDockerRegistries(cfg.ProjectID)
	if err != nil {
		log.Printf("GetListDockerRegistries error: %v", err)
	} else {
		log.Printf("GetListDockerRegistries success: found %d registries", len(registries))
		log.Printf("Registries: %+v", registries)
	}
}

func createDockerRegistry(cfg *config.Config, name string, isPublic bool) {
	dr := cloudru.NewArtifactRegistryApplication(cfg)

	log.Printf("Testing CreateDockerRegistry with name: %s, isPublic: %v...", name, isPublic)
	registry, err := dr.CreateDockerRegistry(
		cfg.ProjectID,
		name,
		isPublic,
	)
	if err != nil {
		log.Printf("CreateDockerRegistry error: %v", err)
	} else {
		log.Printf("CreateDockerRegistry success: %+v", registry)
	}
}
