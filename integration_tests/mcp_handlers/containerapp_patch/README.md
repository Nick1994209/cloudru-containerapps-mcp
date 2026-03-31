# ContainerApp Patch Integration Tests

This directory contains integration tests for the ContainerApp patch functionality, which tests the `internal/presentation/handlers/containerapp_patch.go` handler.

## Overview

These integration tests verify that the patch operation works correctly by:
1. Creating a test container app
2. Applying various patch operations to update different fields
3. Verifying that the changes are applied correctly
4. Cleaning up the test container app

## Test Coverage

The tests cover the following scenarios:

### 1. Basic Fields Patch
- Updates container port
- Updates container image
- Updates description
- Updates publicly accessible flag
- Updates protocol

### 2. Container-Specific Fields Patch
- Updates CPU allocation
- Updates timeout
- Updates idle timeout

### 3. Scaling Fields Patch
- Updates minimum instance count
- Updates maximum instance count

### 4. Auto-Deployments Patch
- Enables/disables auto-deployments
- Updates auto-deployments pattern

### 5. Environment Variables Patch
- Adds multiple environment variables
- Verifies environment variable values

### 6. Command and Args Patch
- Updates container command
- Updates container arguments

### 7. Multiple Fields Patch
- Updates multiple fields simultaneously
- Verifies all changes are applied correctly

### 8. Partial Update Patch
- Updates only specific fields
- Verifies other fields remain unchanged

## Prerequisites

Before running these tests, ensure you have:

1. A valid Cloud.ru account with appropriate permissions
2. Environment variables configured in a `.env` file:
   - `PROJECT_ID`: Your Cloud.ru project ID
   - `KEY_ID`: Your Cloud.ru API key ID
   - `KEY_SECRET`: Your Cloud.ru API key secret
   - `REGISTRY_NAME`: Your Docker registry name
   - `REGISTRY_DOMAIN`: Your Docker registry domain
   - `REPOSITORY_NAME`: Your Docker repository name

## Running the Tests

### From the integration test directory:

```bash
cd integration_tests/mcp_handlers/containerapp_patch
go run containerapp_patch_test.go
```

### From the project root:

```bash
go run integration_tests/mcp_handlers/containerapp_patch/containerapp_patch_test.go
```

## Test Output

The tests will output detailed logs showing:
- Test container app creation
- Each patch operation
- Verification of changes
- Cleanup operations

Example output:
```
=== Creating test container app ===
Test container app created: {...}
Waiting for container app to be ready...
=== Running ContainerApp Patch Integration Tests ===

--- Test: Patch Basic Fields ---
Patch operation completed: {...}
✓ Port updated correctly: 8081
✓ Description updated correctly: Updated description for basic fields test
✓ Publicly accessible updated correctly: false
✓ Protocol updated correctly: http_2
✓ Basic fields patch test passed

--- Test: Patch Container Specific Fields ---
...

=== Cleaning up test container app ===
Test container app deleted successfully
=== All ContainerApp Patch Integration Tests Completed ===
```

## Cleanup

The tests automatically clean up the test container app after completion. However, if a test fails or is interrupted, you may need to manually delete the test container app:

```bash
# Using the MCP server or Cloud.ru CLI
# Delete the container app named: test-patch-{project_id_prefix}
```

## Troubleshooting

### Test fails with authentication error
- Verify your `.env` file contains valid credentials
- Check that your API keys have the necessary permissions

### Test fails with timeout
- Increase the sleep time between operations
- Check your network connection to Cloud.ru

### Container app creation fails
- Verify you have sufficient quota in your Cloud.ru project
- Check that the container image is accessible

## Notes

- The tests use a unique container app name based on your project ID to avoid conflicts
- Each test waits for the container app to be ready before proceeding
- The tests are designed to be run sequentially, not in parallel
- All tests share the same container app, applying patches incrementally

## Related Files

- Handler: `internal/presentation/handlers/containerapp_patch.go`
- Domain types: `internal/domain/types.go`
- Application layer: `internal/application/cloudru/containerapps.go`
