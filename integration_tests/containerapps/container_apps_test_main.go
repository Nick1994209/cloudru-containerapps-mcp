package main

import (
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

func main() {
	cfg := config.LoadConfig()
	containerName := "testme" + cfg.ProjectID

	getListContainerApps(cfg)
	createContainerApp(
		cfg,
		containerName,
		"quickstart.cr.cloud.ru/react-helloworld@sha256:a1a1e0a11299668c5f05a299f74b3943236ca3390a6fda64e98cc2498064c266",
		8080,
	)
	getContainerApp(cfg, containerName)
	deleteContainerApp(cfg, containerName)

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

func createContainerApp(cfg *config.Config, name, image string, port int) {
	ca := cloudru.NewContainerAppsApplication(cfg)

	// Test CreateContainerApp
	request := domain.CreateContainerAppRequest{
		ProjectID:              cfg.ProjectID,
		ContainerAppName:       name,
		ContainerAppPort:       port,
		ContainerAppImage:      image,
		AutoDeploymentsEnabled: false,    // autoDeploymentsEnabled
		AutoDeploymentsPattern: "latest", // autoDeploymentsPattern
		Privileged:             false,    // privileged
		IdleTimeout:            "600s",   // idleTimeout
		Timeout:                "60s",    // timeout
		CPU:                    "0.1",    // cpu
	}

	containerApp, err := ca.CreateContainerApp(request)
	log.Print(containerApp)
	if err != nil {
		log.Fatalf("CreateContainerApp error: %v", err.Error())
	}
}
