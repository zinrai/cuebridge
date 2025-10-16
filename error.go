package cuebridge

import (
	"strings"

	"cuelang.org/go/cue/errors"
)

// extractValidationErrors extracts structured error information from CUE errors
func extractValidationErrors(err error) []ValidationError {
	cueErrors := errors.Errors(err)
	if len(cueErrors) == 0 {
		return []ValidationError{
			{Line: 0, Column: 0, Path: "", Message: err.Error()},
		}
	}

	var validationErrors []ValidationError
	for _, e := range cueErrors {
		ve := extractSingleError(e)
		validationErrors = append(validationErrors, ve)
	}

	return validationErrors
}

// extractSingleError extracts information from a single CUE error
func extractSingleError(e errors.Error) ValidationError {
	return ValidationError{
		Line:    extractLineNumber(e),
		Column:  extractColumnNumber(e),
		Path:    extractFieldPath(e),
		Message: e.Error(),
	}
}

// extractLineNumber extracts the line number from error positions
func extractLineNumber(e errors.Error) int {
	positions := errors.Positions(e)
	for _, pos := range positions {
		if line := pos.Line(); line > 0 {
			return line
		}
	}
	return 0
}

// extractColumnNumber extracts the column number from error positions
func extractColumnNumber(e errors.Error) int {
	positions := errors.Positions(e)
	for _, pos := range positions {
		if col := pos.Column(); col > 0 {
			return col
		}
	}
	return 0
}

// extractFieldPath extracts and formats the field path from error
func extractFieldPath(e errors.Error) string {
	path := e.Path()
	if len(path) == 0 {
		return ""
	}
	return formatPath(path)
}

// formatPath converts CUE path to a string representation
func formatPath(path []string) string {
	var parts []string

	for _, p := range path {
		if !isValidPathElement(p) {
			continue
		}
		p = strings.Trim(p, `"`)
		parts = append(parts, p)
	}

	return strings.Join(parts, ".")
}

// isValidPathElement checks if a path element should be included
func isValidPathElement(p string) bool {
	return p != "" &&
		!strings.HasPrefix(p, "[") &&
		p != "#Config"
}
