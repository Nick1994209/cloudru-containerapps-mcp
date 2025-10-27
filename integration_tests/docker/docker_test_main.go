package main

import (
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

func main() {
	cfg := config.LoadConfig()

	// Test Docker login
	// testDockerLogin(cfg, "korolkov-mcp-lichniy")

	// Test Docker build and push commands (without executing)
	// testShowBuildAndPushCommands(cfg)

	// Test Docker build and push (requires actual Docker setup)
	// testBuildAndPush(cfg)

	// Test getting registry images
	testGetRegistryImages(cfg, "korolkov-mcp-lichniy")
}

func testDockerLogin(cfg *config.Config, registryName string) {
	dockerApp := application.NewDockerApplication(cfg)

	log.Printf("Testing Docker login with registry: %s...", registryName)
	loginTarget, err := dockerApp.Login(registryName)
	if err != nil {
		log.Printf("Docker login error: %v", err)
	} else {
		log.Printf("Docker login success: logged in to %s", loginTarget)
	}
}

func testShowBuildAndPushCommands(cfg *config.Config) {
	dockerApp := application.NewDockerApplication(cfg)

	image := domain.DockerImage{
		RegistryName:     "nvkorolkov",
		RepositoryName:   "test-repo",
		ImageVersion:     "latest",
		DockerfilePath:   "Dockerfile",
		DockerfileTarget: "",
		DockerfileFolder: ".",
	}

	log.Println("Testing ShowBuildAndPushCommands...")
	buildCmd, pushCmd, err := dockerApp.ShowBuildAndPushCommands(image)
	if err != nil {
		log.Printf("ShowBuildAndPushCommands error: %v", err)
	} else {
		log.Printf("ShowBuildAndPushCommands success:")
		log.Printf("Build command: %s", buildCmd)
		log.Printf("Push command: %s", pushCmd)
	}
}

func testBuildAndPush(cfg *config.Config) {
	dockerApp := application.NewDockerApplication(cfg)

	image := domain.DockerImage{
		RegistryName:     "nvkorolkov",
		RepositoryName:   "test-repo",
		ImageVersion:     "latest",
		DockerfilePath:   "Dockerfile",
		DockerfileTarget: "",
		DockerfileFolder: ".",
	}

	log.Println("Testing BuildAndPush...")
	imageTag, err := dockerApp.BuildAndPush(image)
	if err != nil {
		log.Printf("BuildAndPush error: %v", err)
	} else {
		log.Printf("BuildAndPush success: pushed image %s", imageTag)
	}
}

func testGetRegistryImages(cfg *config.Config, registryName string) {
	dockerApp := application.NewDockerApplication(cfg)

	log.Printf("Testing GetRegistryImages with registry: %s...", registryName)
	images, err := dockerApp.GetRegistryImages(registryName)
	if err != nil {
		log.Printf("GetRegistryImages error: %v", err)
	} else {
		log.Printf("GetRegistryImages success: found %d images", len(images))
		for _, img := range images {
			log.Printf("Image: %+v", img)
		}
	}
}
