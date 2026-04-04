package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// Helper function to format request body as JSON string
func formatRequestBody(body interface{}) string {
	if body == nil {
		return "nil"
	}
	jsonBytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to marshal request body: %v", err)
	}
	return string(jsonBytes)
}

// Helper function to log error with request body
func logErrorWithRequestBody(message string, err error, requestBody interface{}) {
	log.Printf("%s\nError: %v\nRequest Body:\n%s", message, err, formatRequestBody(requestBody))
}

// Helper functions for creating pointer values (similar to containerapp_patch)
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// validateOperation validates an operation response
func validateOperation(operation *domain.Operation) bool {
	// Check if operation is nil
	if operation == nil {
		log.Printf("Validation failed: operation is nil")
		return false
	}

	// Check if ResourceID is not empty
	if operation.ResourceID == "" {
		log.Printf("Validation failed: ResourceID is empty")
		return false
	}

	// Check if ResourceName is not empty
	if operation.ResourceName == "" {
		log.Printf("Validation failed: ResourceName is empty")
		return false
	}

	// Check if Description is not empty
	if operation.Description == "" {
		log.Printf("Validation failed: Description is empty")
		return false
	}

	// Log operation details for verification
	log.Printf("Operation validated - ResourceID: %s, ResourceName: %s, Description: %s, Done: %v",
		operation.ResourceID, operation.ResourceName, operation.Description, operation.Done)

	return true
}
