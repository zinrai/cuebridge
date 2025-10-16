package cuebridge

import (
	"fmt"
	"strings"
)

// FormatResults formats validation results into a human-readable string.
//
// Output format:
//
//	config.yaml: ok
//	config.json: FAIL
//	  line 5, field "replicas": value 0 does not satisfy constraint >=1
func FormatResults(results []ValidationResult) string {
	var output strings.Builder

	for _, result := range results {
		formatSingleResult(&output, result)
	}

	return output.String()
}

// formatSingleResult formats a single validation result
func formatSingleResult(output *strings.Builder, result ValidationResult) {
	if result.Valid {
		fmt.Fprintf(output, "%s: ok\n", result.Name)
		return
	}

	fmt.Fprintf(output, "FAIL: %s\n", result.Name)
	for _, err := range result.Errors {
		formatError(output, err)
	}
}

// formatError formats a single validation error
func formatError(output *strings.Builder, err ValidationError) {
	switch {
	case err.Line > 0 && err.Path != "":
		fmt.Fprintf(output, "  line %d, field \"%s\": %s\n",
			err.Line, err.Path, err.Message)
	case err.Line > 0:
		fmt.Fprintf(output, "  line %d: %s\n",
			err.Line, err.Message)
	case err.Path != "":
		fmt.Fprintf(output, "  field \"%s\": %s\n",
			err.Path, err.Message)
	default:
		fmt.Fprintf(output, "  %s\n", err.Message)
	}
}
