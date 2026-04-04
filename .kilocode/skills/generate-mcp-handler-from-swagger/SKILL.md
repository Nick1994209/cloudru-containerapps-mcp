---
name: generate-mcp-handler-from-swagger
description: Generate MCP handler methods based on Swagger/OpenAPI specifications for Cloud.ru Container Apps
---

# Generate MCP Handler from Swagger

## When to use this skill
Use this skill when you need to generate new MCP handler methods based on Swagger/OpenAPI specifications for Cloud.ru Container Apps service. This skill helps automate the creation of consistent handler code that follows the established patterns in the project.

## Prerequisites
- Valid OpenAPI/Swagger specification for the Cloud.ru Container Apps API endpoint you want to implement
- Understanding of the existing MCP handler patterns in the project

## Required Parameters
The following parameters are mandatory:
- `openapi_spec` - The OpenAPI/Swagger specification (JSON/YAML format) for the endpoint
- `handler_name` - Name for the new handler (e.g., "containerapp_create")
- `method_name` - The HTTP method name from the spec (e.g., "CreateContainerApp")

## Process Overview
1. Validate the OpenAPI specification
2. Extract relevant information from the spec:
   - Operation ID
   - Parameters and their types
   - Request/Response schemas
   - Required vs optional fields
3. Generate the handler code following the established patterns
4. Create appropriate domain types if needed
5. Register the tool in the server

## How to validate OpenAPI specification
1. Ensure the specification is valid JSON or YAML
2. Verify it contains the operation you want to implement
3. Check that all required parameters and schemas are defined
4. Confirm the specification matches the Cloud.ru Container Apps API structure

## Generated Handler Structure
The generated handler will include:

### Registration Function
```go
func (s *MCPServer) Register[MethodName]Tool(mcpServer *server.MCPServer) {
    // Tool registration with parameters
}
```

### Parameter Processing
- Extraction of required and optional parameters
- Type conversion (string to int, bool, etc.)
- Default value handling
- Error checking for required parameters

### Service Call
- Creation of domain request structs
- Calling the appropriate service method
- Handling responses and errors

### Response Formatting
- JSON marshaling of results
- Proper error handling and messaging

## Examples

### Sample Input OpenAPI Spec (YAML)
```yaml
paths:
  /projects/{projectId}/containerApps:
    post:
      operationId: CreateContainerApp
      summary: Create a new Container App
      parameters:
        - name: projectId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateContainerAppRequest'
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Operation'
```

### Generated Handler (simplified)
```go
func (s *MCPServer) RegisterCreateContainerAppTool(mcpServer *server.MCPServer) {
    toolOptions := s.getMCPFieldsOptions(
        "Create a new Container App in Cloud.ru",
        "project_id",
        "containerapp_name",
        "containerapp_port",
        // ... other fields
    )
    createContainerAppTool := mcp.NewTool("cloudru_create_containerapp", toolOptions...)

    mcpServer.AddTool(createContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Extract parameters
        projectID, err := s.getMCPFieldValue("project_id", request)
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        // ... more parameter extraction

        // Create request struct
        createRequest := domain.CreateContainerAppRequest{
            ProjectID: projectID,
            // ... other fields
        }

        // Call service
        operation, err := s.containerAppsService.CreateContainerApp(createRequest)
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        // Format response
        result, err := json.MarshalIndent(operation, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
        }

        return mcp.NewToolResultText(fmt.Sprintf("Successfully created Container App: %s\n%s", containerAppName, string(result))), nil
    })
}
```

## Important Notes
1. The OpenAPI specification is REQUIRED - without it, the generation cannot proceed
2. Generated code should be reviewed and adjusted as needed
3. Ensure parameter names match the existing conventions in mappedFields
4. Follow the existing patterns for error handling and response formatting
5. Update the RegisterAllTools function to register the new handler