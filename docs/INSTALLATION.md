## Prerequisites

- Go 1.22 or later
- Docker installed and configured
- Access to Cloud.ru Container Apps service
- Service account credentials with appropriate permissions

## Installation

There are several ways to install and use this MCP:

### Method 1: Using go install (Recommended for released versions)

If you have Go installed on your system and want to install a released version from GitHub, you can use:

```bash
go install github.com/Nick1994209/cloudru-containerapps-mcp/cmd/cloudru-containerapps-mcp@latest
```

This will download, compile, and install the binary to your `$GOPATH/bin` directory.

Note: This method works for released versions that include the proper directory structure.
For local development or if the remote repository doesn't have the cmd directory yet,
use the build from source method below.

### Method 2: Building from source

1. Clone this repository
2. Run `go build -o cloudru-containerapps-mcp cmd/cloudru-containerapps-mcp/main.go` to build the binary
3. Make sure Docker is installed and running on your system

## Making Go Binaries Available in Your PATH

To use Go-installed binaries from anywhere in your system, you need to ensure your `$GOPATH/bin` directory is in your system PATH.

### Finding Your GOPATH

First, check your GOPATH:

```bash
go env GOPATH
```

By default, this is usually `$HOME/go`.

### Adding GOPATH/bin to Your PATH

#### For Bash Users

Add this line to your `~/.bashrc` or `~/.bash_profile`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell configuration:

```bash
source ~/.bashrc
# or
source ~/.bash_profile
```

#### For Zsh Users

Add this line to your `~/.zshrc`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell configuration:

```bash
source ~/.zshrc
```

#### For Fish Users

Add this line to your `~/.config/fish/config.fish`:

```fish
set -gx PATH $PATH (go env GOPATH)/bin
```

Then reload your shell configuration:

```bash
source ~/.config/fish/config.fish
```

### Verifying the Installation

After adding GOPATH/bin to your PATH, you can verify that the binary is accessible:

```bash
cloudru-containerapps-mcp
```

Note: Since this is an MCP server that communicates via stdin/stdout, running it directly might not produce visible output. It's meant to be used with MCP-compatible clients like Kilo Code, Roo Code, or Claude.

## After installation

[docs/HOW_ADD_TO_IDE.md](HOW_ADD_TO_IDE.md)
