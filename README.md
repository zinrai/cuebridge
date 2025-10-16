# cuebridge

A Go library that bridges JSON/YAML data with CUE validation logic.

## Overview

`cuebridge` connects JSON and YAML configuration files with [CUE](https://cuelang.org/) schemas for validation.

**What it does:**

- Reads data from files, readers, or byte slices
- Parses JSON and YAML into CUE values
- Evaluates data against CUE schemas
- Extracts detailed error information
- Formats validation results as text

**What it does not do:**

- Define validation rules (delegated to CUE)
- Detect file formats automatically (caller specifies)
- Iterate over multiple files (caller's responsibility)
- Decide exit codes (caller's responsibility)

## Installation

```bash
$ go get github.com/zinrai/cuebridge
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zinrai/cuebridge"
)

func main() {
    // Create validator with CUE schema
    validator, err := cuebridge.NewValidator("schema.cue", "#Config")
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate a YAML file
    result, err := validator.Validate(cuebridge.ValidationInput{
        SourceType: cuebridge.SourceFile,
        FilePath:   "config.yaml",
        Format:     cuebridge.FormatYAML,
        Name:       "config.yaml",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Check result
    if !result.Valid {
        fmt.Println("Validation failed:")
        for _, e := range result.Errors {
            fmt.Printf("  Line %d: %s\n", e.Line, e.Message)
        }
    }
}
```

### Validating Multiple Files

```go
validator, _ := cuebridge.NewValidator("schema.cue", "#Config")

var results []cuebridge.ValidationResult
for _, path := range configPaths {
    // Determine format from extension
    format := cuebridge.FormatYAML
    if strings.HasSuffix(path, ".json") {
        format = cuebridge.FormatJSON
    }
    
    result, err := validator.Validate(cuebridge.ValidationInput{
        SourceType: cuebridge.SourceFile,
        FilePath:   path,
        Format:     format,
        Name:       path,
    })
    if err != nil {
        log.Fatal(err)
    }
    results = append(results, result)
}

// Format and display results
output := cuebridge.FormatResults(results)
fmt.Print(output)
```

### Using Different Definition Names

```go
// Use #ServiceConfig instead of #Config
validator, err := cuebridge.NewValidator("schema.cue", "#ServiceConfig")

// Use #Application
validator, err := cuebridge.NewValidator("schema.cue", "#Application")
```

### Reading from stdin

```go
result, err := validator.Validate(cuebridge.ValidationInput{
    SourceType: cuebridge.SourceReader,
    Reader:     os.Stdin,
    Format:     cuebridge.FormatJSON,
    Name:       "stdin",
})
```

### Using Byte Slices

```go
data := []byte(`{"name": "test", "replicas": 3}`)

result, err := validator.Validate(cuebridge.ValidationInput{
    SourceType: cuebridge.SourceBytes,
    Data:       data,
    Format:     cuebridge.FormatJSON,
    Name:       "config",
})
```

## Example CUE Schema

```cue
package example

#Config: {
    // Application name (required)
    name: string & =~"^[a-z][a-z0-9-]*$"
    
    // Number of replicas (optional, 1-10)
    replicas?: int & >=1 & <=10
    
    // Deployment environment
    environment: "development" | "staging" | "production"
}
```

The definition name (e.g., `#Config`) can be anything. Specify it when creating a validator:

```go
validator, _ := cuebridge.NewValidator("schema.cue", "#Config")
```

## Output Format

Results are formatted as human-readable text:

```
config.yaml: ok
config.json: FAIL
  line 5, field "replicas": value 0 does not satisfy constraint >=1
```

For CI/CD integration, check the `Valid` field and use appropriate exit codes in your tool.

## Supported Formats

**Input formats:**

- JSON (`.json`)
- YAML (`.yaml`, `.yml`)

**Output format:**

- Text (human-readable)

## Real-World Usage

See these projects using cuebridge:

- [ycint](https://github.com/zinrai/ycint) - YAML-only configuration linter
- [integratify](https://github.com/zinrai/integratify) - CI/CD integration tool

## Design Principles

1. **Delegation to CUE**: All validation logic is defined in CUE schemas, not in Go code
2. **Explicit parameters**: Caller specifies format (JSON/YAML) and definition name (e.g., `#Config`)
3. **Three input sources**: Files, io.Reader, or byte slices
4. **Single responsibility**: Only handles data reading, parsing, CUE evaluation, and result formatting
5. **Caller control**: File iteration, format detection, and exit codes are the caller's responsibility
6. **Minimal API**: Three public functions (`NewValidator`, `Validate`, `FormatResults`)

## License

This project is licensed under the [MIT License](./LICENSE).
