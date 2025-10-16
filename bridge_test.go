package cuebridge

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestEndToEnd tests the complete validation flow
func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name      string
		schema    string
		content   string
		filename  string
		format    DataFormat
		wantValid bool
	}{
		{
			name:      "yaml valid",
			schema:    `#Config: {name: string}`,
			content:   `name: "test"`,
			filename:  "config.yaml",
			format:    FormatYAML,
			wantValid: true,
		},
		{
			name:      "json valid",
			schema:    `#Config: {name: string}`,
			content:   `{"name": "test"}`,
			filename:  "config.json",
			format:    FormatJSON,
			wantValid: true,
		},
		{
			name:      "yaml invalid",
			schema:    `#Config: {name: string}`,
			content:   `wrong: "field"`,
			filename:  "config.yaml",
			format:    FormatYAML,
			wantValid: false,
		},
		{
			name:      "json invalid",
			schema:    `#Config: {name: string}`,
			content:   `{"wrong": "field"}`,
			filename:  "config.json",
			format:    FormatJSON,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			schemaPath := filepath.Join(tmpDir, "schema.cue")
			if err := os.WriteFile(schemaPath, []byte(tt.schema), 0644); err != nil {
				t.Fatalf("failed to write schema: %v", err)
			}

			configPath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write config: %v", err)
			}

			validator, err := NewValidator(schemaPath, "#Config")
			if err != nil {
				t.Fatalf("NewValidator failed: %v", err)
			}

			result, err := validator.Validate(ValidationInput{
				SourceType: SourceFile,
				FilePath:   configPath,
				Format:     tt.format,
				Name:       tt.filename,
			})
			if err != nil {
				t.Fatalf("Validate failed: %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

// TestInputSources tests different input sources produce the same result
func TestInputSources(t *testing.T) {
	tmpDir := t.TempDir()

	schema := `#Config: {name: string}`
	schemaPath := filepath.Join(tmpDir, "schema.cue")
	if err := os.WriteFile(schemaPath, []byte(schema), 0644); err != nil {
		t.Fatalf("failed to write schema: %v", err)
	}

	validator, err := NewValidator(schemaPath, "#Config")
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	testData := []byte(`name: "test"`)

	// Test 1: SourceFile
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, testData, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	resultFile, err := validator.Validate(ValidationInput{
		SourceType: SourceFile,
		FilePath:   configPath,
		Format:     FormatYAML,
		Name:       "file",
	})
	if err != nil {
		t.Fatalf("Validate (SourceFile) failed: %v", err)
	}

	// Test 2: SourceReader
	reader := bytes.NewReader(testData)
	resultReader, err := validator.Validate(ValidationInput{
		SourceType: SourceReader,
		Reader:     reader,
		Format:     FormatYAML,
		Name:       "reader",
	})
	if err != nil {
		t.Fatalf("Validate (SourceReader) failed: %v", err)
	}

	// Test 3: SourceBytes
	resultBytes, err := validator.Validate(ValidationInput{
		SourceType: SourceBytes,
		Data:       testData,
		Format:     FormatYAML,
		Name:       "bytes",
	})
	if err != nil {
		t.Fatalf("Validate (SourceBytes) failed: %v", err)
	}

	// All should produce the same result
	if !resultFile.Valid || !resultReader.Valid || !resultBytes.Valid {
		t.Error("all input sources should produce valid result")
	}
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		setup     func() (string, ValidationInput)
		wantError bool
	}{
		{
			name: "non-existent file",
			setup: func() (string, ValidationInput) {
				schema := `#Config: {name: string}`
				schemaPath := filepath.Join(tmpDir, "schema.cue")
				os.WriteFile(schemaPath, []byte(schema), 0644)
				return schemaPath, ValidationInput{
					SourceType: SourceFile,
					FilePath:   "/nonexistent/file.yaml",
					Format:     FormatYAML,
					Name:       "missing",
				}
			},
			wantError: true,
		},
		{
			name: "invalid schema",
			setup: func() (string, ValidationInput) {
				schemaPath := filepath.Join(tmpDir, "invalid.cue")
				os.WriteFile(schemaPath, []byte(`broken syntax`), 0644)
				return schemaPath, ValidationInput{}
			},
			wantError: true,
		},
		{
			name: "schema without Config",
			setup: func() (string, ValidationInput) {
				schemaPath := filepath.Join(tmpDir, "no-config.cue")
				os.WriteFile(schemaPath, []byte(`package test`), 0644)
				return schemaPath, ValidationInput{}
			},
			wantError: true,
		},
		{
			name: "invalid json syntax",
			setup: func() (string, ValidationInput) {
				schema := `#Config: {name: string}`
				schemaPath := filepath.Join(tmpDir, "schema2.cue")
				os.WriteFile(schemaPath, []byte(schema), 0644)

				configPath := filepath.Join(tmpDir, "invalid.json")
				os.WriteFile(configPath, []byte(`{"name": "test"`), 0644) // missing closing brace

				return schemaPath, ValidationInput{
					SourceType: SourceFile,
					FilePath:   configPath,
					Format:     FormatJSON,
					Name:       "invalid.json",
				}
			},
			wantError: false, // Parse error returns as ValidationResult, not error
		},
		{
			name: "invalid yaml syntax",
			setup: func() (string, ValidationInput) {
				schema := `#Config: {name: string}`
				schemaPath := filepath.Join(tmpDir, "schema3.cue")
				os.WriteFile(schemaPath, []byte(schema), 0644)

				configPath := filepath.Join(tmpDir, "invalid.yaml")
				os.WriteFile(configPath, []byte("name:\n  - bad indent"), 0644)

				return schemaPath, ValidationInput{
					SourceType: SourceFile,
					FilePath:   configPath,
					Format:     FormatYAML,
					Name:       "invalid.yaml",
				}
			},
			wantError: false, // Parse error returns as ValidationResult, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaPath, input := tt.setup()

			validator, err := NewValidator(schemaPath, "#Config")
			if tt.name == "invalid schema" || tt.name == "schema without Config" {
				if err == nil {
					t.Error("expected error from NewValidator")
				}
				return
			}
			if err != nil && !tt.wantError {
				t.Fatalf("NewValidator failed: %v", err)
			}

			_, err = validator.Validate(input)
			if tt.wantError && err == nil {
				t.Error("expected error from Validate")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error from Validate: %v", err)
			}
		})
	}
}
