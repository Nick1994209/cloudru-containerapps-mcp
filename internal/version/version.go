package version

import (
	"fmt"
	"runtime"
)

// Version information
var (
	// Version is the application version
	Version = "0.0.1"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// BuildDate is the build date
	BuildDate = "unknown"

	// GoVersion is the Go version used to build
	GoVersion = runtime.Version()
)

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	return fmt.Sprintf("Cloud.ru Container Apps MCP\nVersion: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s",
		Version, GitCommit, BuildDate, GoVersion)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}
