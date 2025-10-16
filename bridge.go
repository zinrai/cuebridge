// Package cuebridge provides a bridge between JSON/YAML data and CUE validation logic.
//
// cuebridge is a library that connects JSON/YAML configuration files with CUE schemas
// for validation. It handles reading data from various sources (files, readers, bytes),
// parsing them into CUE values, evaluating against schemas, and formatting results.
//
// Basic usage:
//
//	validator, err := cuebridge.NewValidator("schema.cue", "#Config")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	result, err := validator.Validate(cuebridge.ValidationInput{
//	    SourceType: cuebridge.SourceFile,
//	    FilePath:   "config.yaml",
//	    Format:     cuebridge.FormatYAML,
//	    Name:       "config.yaml",
//	})
//
//	if !result.Valid {
//	    for _, e := range result.Errors {
//	        fmt.Printf("Error at line %d: %s\n", e.Line, e.Message)
//	    }
//	}
package cuebridge

import (
	"io"

	"cuelang.org/go/cue"
)

// InputSourceType represents the type of input source
type InputSourceType int

const (
	// SourceFile reads data from a file path
	SourceFile InputSourceType = iota
	// SourceReader reads data from an io.Reader
	SourceReader
	// SourceBytes uses data directly from a byte slice
	SourceBytes
)

// DataFormat represents the format of input data
type DataFormat int

const (
	// FormatJSON represents JSON format
	FormatJSON DataFormat = iota
	// FormatYAML represents YAML format
	FormatYAML
)

// Validator validates data against a CUE schema.
// Create a Validator with NewValidator and reuse it for multiple validations.
type Validator struct {
	schemaPath     string
	definitionName string
	ctx            *cue.Context
	compiledSchema cue.Value
}

// ValidationInput specifies the input data to validate.
type ValidationInput struct {
	// SourceType specifies where to read data from
	SourceType InputSourceType
	// Name is an identifier for this input (used in error messages)
	Name string
	// FilePath is the path to the file (when SourceType is SourceFile)
	FilePath string
	// Reader is the io.Reader to read from (when SourceType is SourceReader)
	Reader io.Reader
	// Data is the byte slice to use directly (when SourceType is SourceBytes)
	Data []byte
	// Format specifies the data format (FormatJSON or FormatYAML)
	Format DataFormat
}

// ValidationResult contains the result of validating a single input.
type ValidationResult struct {
	// Name is the identifier of the validated input
	Name string
	// Valid is true if validation succeeded
	Valid bool
	// Errors contains validation errors (empty if Valid is true)
	Errors []ValidationError
}

// ValidationError represents a single validation error.
type ValidationError struct {
	// Line is the line number in the source (0 if unknown)
	Line int
	// Column is the column number in the source (0 if unknown)
	Column int
	// Path is the field path (e.g., "spec.replicas")
	Path string
	// Message is the error message
	Message string
}

// NewValidator creates a new Validator by loading and compiling a CUE schema file.
// The definitionName parameter specifies which definition to use for validation
// (e.g., "#Config", "#ServiceConfig").
//
// Returns an error if:
//   - The schema file cannot be read
//   - The schema has CUE syntax errors
//   - The schema does not define the specified definition
func NewValidator(schemaPath string, definitionName string) (*Validator, error) {
	return newValidator(schemaPath, definitionName)
}

// Validate validates a single input against the schema.
//
// Returns ValidationResult with Valid=false if validation fails.
// Returns an error only if the validation process itself fails
// (e.g., file cannot be read, unsupported format).
func (v *Validator) Validate(input ValidationInput) (ValidationResult, error) {
	return v.validate(input)
}
