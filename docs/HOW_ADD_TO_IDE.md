# How to Add Cloud.ru Container Apps MCP to IDEs and AI Coding Assistants

This guide explains how to integrate the Cloud.ru Container Apps MCP with popular IDEs and AI coding assistants.

## Prerequisites

Before adding the MCP to your IDE, ensure you have:

1. Installed  [INSTALLATION.md](INSTALLATION.md)
2. Set up your Cloud.ru credentials as environment variables:
   Alternatively, you can create a `.env` file in the project root directory with these variables.
   
   **Required environment variables:**
   - `CLOUDRU_KEY_ID` - Your service account key ID (required)
   - `CLOUDRU_KEY_SECRET` - Your service account key secret (required)
   
   https://cloud.ru/docs/console_api/ug/topics/quickstart

   More about [docs/ENVIRONMENT_VARIABLES.md](ENVIRONMENT_VARIABLES.md)   

## Kilo Code

### Adding the MCP Server

1. Open Kilo Code settings
2. Navigate to "Tools" → "MCP Servers"
3. Click "Add New Server"
4. Configure with:
   - Name: `Cloud.ru Container Apps`
   - Command: `cloudru-containerapps-mcp`

### Example Configuration JSON

Here's a complete example of how to configure this MCP server in Kilo Code:

```json
{
  "mcpServers": {
    "cloudru-containerapps-mcp": {
      "command": "cloudru-containerapps-mcp",
      "args": [],
      "env": {
        "CLOUDRU_KEY_ID": "(REQUIRED FIELD) your-service-account-key-id",
        "CLOUDRU_KEY_SECRET": "(REQUIRED FIELD) your-service-account-key-secret",
        "CLOUDRU_REGISTRY_NAME": "your-registry-name",
        "CLOUDRU_PROJECT_ID": "your-project-id",
        "CLOUDRU_REPOSITORY_NAME": "your-repository-name",
        "CLOUDRU_DOCKERFILE": "Dockerfile",
        "CLOUDRU_DOCKERFILE_TARGET": "-",
        "CLOUDRU_DOCKERFILE_FOLDER": "."
      }
    }
  }
}
```

4. Click "Save and Connect"
5. The tools should now be available in Kilo Code with access to your environment variables

## Roo Code

1. Open Roo Code preferences
2. Go to "Extensions" → "MCP Integration"
3. Click "Add MCP Server"
4. Fill in the details:
   - Server Name: `Cloud.ru Container Apps`
   - Executable Path: Full path to `cloudru-containerapps-mcp` (or just `cloudru-containerapps-mcp` if installed via `go install` and GOPATH/bin is in your PATH)
   - Arguments: Leave empty
5. Configure environment variables in the server settings if supported
6. Click "Test Connection" to verify
7. Enable the server and restart Roo Code if needed

## Claude (Anthropic Console)

1. Access the Anthropic Console
2. Navigate to "Tools" → "Custom MCPs"
3. Click "Register New MCP"
4. Provide the following configuration:
   - Tool Name: `cloudru-containerapps`
   - Execution Method: `subprocess`
   - Command: `["cloudru-containerapps-mcp"]`
   - Working Directory: Project path (only needed if using the local binary)
   - Environment Variables: Add as key-value pairs in the configuration interface
5. Save the configuration
6. The tools will be available in Claude prompts

## Cursor

1. Open Cursor settings
2. Go to "Extensions" → "MCP Tools"
3. Click "Add External Tool"
4. Enter configuration:
   - Tool Identifier: `cloudru-containerapps-mcp`
   - Executable: Path to the `cloudru-containerapps-mcp` binary
   - Environment Variables: Set in the tool configuration if supported
   - Auto-start: Enabled
5. Restart Cursor to load the tool
6. Verify integration by opening the tools panel

## General Usage Notes

Once integrated, you can use the following tools:

1. `cloudru_containerapps_description()` - Get usage instructions
2. `cloudru_docker_login()` - Authenticate with Cloud.ru registry
3. `cloudru_docker_build_and_push()` - Build and push Docker images (supports dockerfile_path, dockerfile_target, and dockerfile_folder parameters)

## Example prompts
[docs/EXAMPLE_PROMPTS.md](EXAMPLE_PROMPTS.md)

## Troubleshooting

If the MCP doesn't appear in your IDE:

1. Verify the binary is executable: `chmod +x cloudru-containerapps-mcp`
2. Test the binary directly: `echo '{"method":"initialize","id":1}' | cloudru-containerapps-mcp`
3. Check that all required environment variables are set
4. Ensure your IDE has permission to execute subprocesses
5. Consult your IDE's documentation for MCP/tool integration specifics

## Security Considerations

- The MCP requires access to Docker daemon
- Authentication credentials are handled through environment variables
- Ensure your IDE/assistant has appropriate permissions for container operations
- Never commit sensitive credentials to version control
