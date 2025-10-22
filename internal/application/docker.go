package application

import (
	"fmt"
	"os/exec"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// DockerApplication implements the DockerService interface using actual Docker commands
type DockerApplication struct{}

// NewDockerApplication creates a new DockerApplication
func NewDockerApplication() domain.DockerService {
	return &DockerApplication{}
}

// Login logs into the Cloud.ru Docker registry using Docker CLI
func (d *DockerApplication) Login(registryName string, credentials domain.Credentials) (string, error) {
	loginTarget := fmt.Sprintf("%s.cr.cloud.ru", registryName)
	cmd := exec.Command("docker", "login", loginTarget, "-u", credentials.KeyID, "--password-stdin")

	// Create a pipe to send the password to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Execute the command
	go func() {
		defer stdin.Close()
		fmt.Fprint(stdin, credentials.KeySecret)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker login to %s failed: %w\nOutput: %s\n\nPlease ensure:\n1. The registry exists in Cloud.ru Evolution Artifact Registry\n2. You have created a registry and obtained access keys\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", loginTarget, err, string(output))
	}

	return loginTarget, nil
}

// BuildAndPush builds and pushes a Docker image to Cloud.ru Artifact Registry
func (d *DockerApplication) BuildAndPush(image domain.DockerImage, credentials domain.Credentials) (string, error) {
	imageTag := fmt.Sprintf("%s.cr.cloud.ru/%s:%s", image.RegistryName, image.RepositoryName, image.ImageVersion)

	// Build the Docker image
	var buildCmd *exec.Cmd
	// Set build context folder, default to current directory if not specified
	buildContext := "."
	if image.DockerfileFolder != "" && image.DockerfileFolder != "." {
		buildContext = image.DockerfileFolder
	}

	if image.DockerfileTarget != "" && image.DockerfileTarget != "-" {
		buildCmd = exec.Command("docker", "build", "--platform", "linux/amd64", "-t", imageTag, "--target", image.DockerfileTarget, "-f", image.DockerfilePath, buildContext)
	} else {
		buildCmd = exec.Command("docker", "build", "--platform", "linux/amd64", "-t", imageTag, "-f", image.DockerfilePath, buildContext)
	}
	buildOutput, buildErr := buildCmd.CombinedOutput()

	// Always include build output in the response for visibility
	if len(buildOutput) > 0 {
		fmt.Printf("Docker build output:\n%s\n", string(buildOutput))
	}

	if buildErr != nil {
		return "", fmt.Errorf("failed to build Docker image %s: %w\nOutput: %s", imageTag, buildErr, string(buildOutput))
	}

	// Push the Docker image
	pushCmd := exec.Command("docker", "push", "--platform", "linux/amd64", imageTag)
	pushOutput, pushErr := pushCmd.CombinedOutput()

	// Always include push output in the response for visibility
	if len(pushOutput) > 0 {
		fmt.Printf("Docker push output:\n%s\n", string(pushOutput))
	}

	if pushErr != nil {
		// If push fails, try to re-login if credentials are available
		if (credentials.KeyID != "") || (credentials.KeySecret != "") {
			fmt.Println("Attempting to re-login to Docker registry...")
			_, loginErr := d.Login(image.RegistryName, credentials)
			if loginErr != nil {
				return "", fmt.Errorf("docker push failed and re-login unsuccessful: %w\nOutput: %s\n\nTo resolve this issue:\n1. Set KEY_ID and KEY_SECRET environment variables\n2. Or run the cloudru_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", loginErr, string(pushOutput))
			}

			// Retry push after re-login
			fmt.Println("Retrying Docker push after re-login...")
			pushCmd = exec.Command("docker", "push", imageTag)
			pushOutput, pushErr = pushCmd.CombinedOutput()

			// Always include push output in the response for visibility
			if len(pushOutput) > 0 {
				fmt.Printf("Docker push retry output:\n%s\n", string(pushOutput))
			}

			if pushErr != nil {
				return "", fmt.Errorf("docker push still failed after re-login: %w\nOutput: %s", pushErr, string(pushOutput))
			}
		} else {
			return "", fmt.Errorf("docker push failed: %w\nOutput: %s\n\nTo resolve this issue:\n1. Set KEY_ID and KEY_SECRET environment variables\n2. Or run the cloudru_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", pushErr, string(pushOutput))
		}
	}

	return imageTag, nil
}
