package cuebridge

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

// newValidator creates a new Validator by loading and compiling a CUE schema
func newValidator(schemaPath string, definitionName string) (*Validator, error) {
	// Read schema file
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	// Create CUE context
	ctx := cuecontext.New()

	// Compile schema
	schema := ctx.CompileBytes(schemaData, cue.Filename(schemaPath))
	if schema.Err() != nil {
		return nil, fmt.Errorf("compiling schema: %w", schema.Err())
	}

	// Verify definition exists
	configDef := schema.LookupPath(cue.ParsePath(definitionName))
	if !configDef.Exists() {
		return nil, fmt.Errorf("schema does not define %s", definitionName)
	}

	return &Validator{
		schemaPath:     schemaPath,
		definitionName: definitionName,
		ctx:            ctx,
		compiledSchema: schema,
	}, nil
}

// validate validates a single input against the schema
func (v *Validator) validate(input ValidationInput) (ValidationResult, error) {
	// Read input data
	data, err := readInput(input)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("reading input: %w", err)
	}

	// Parse data into CUE value
	parsedData, err := parseData(v.ctx, data, input.Format, input.Name)
	if err != nil {
		return createErrorResult(input.Name, fmt.Sprintf("failed to parse: %v", err)), nil
	}

	// Check for parse errors
	if parsedData.Err() != nil {
		return createValidationErrorResult(input.Name, parsedData.Err()), nil
	}

	// Get definition from schema
	configDef := v.compiledSchema.LookupPath(cue.ParsePath(v.definitionName))
	if !configDef.Exists() {
		return ValidationResult{}, fmt.Errorf("schema does not define %s", v.definitionName)
	}

	// Unify data with schema
	unified := configDef.Unify(parsedData)

	// Validate
	err = unified.Validate(cue.Concrete(true))
	if err != nil {
		return createValidationErrorResult(input.Name, err), nil
	}

	// Success
	return ValidationResult{
		Name:   input.Name,
		Valid:  true,
		Errors: []ValidationError{},
	}, nil
}

// createErrorResult creates a result with a single error message
func createErrorResult(name string, message string) ValidationResult {
	return ValidationResult{
		Name:  name,
		Valid: false,
		Errors: []ValidationError{
			{Line: 0, Column: 0, Path: "", Message: message},
		},
	}
}

// createValidationErrorResult creates a result with extracted validation errors
func createValidationErrorResult(name string, err error) ValidationResult {
	return ValidationResult{
		Name:   name,
		Valid:  false,
		Errors: extractValidationErrors(err),
	}
}
