package cuebridge

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/yaml"
)

// parseData parses data into a CUE value based on format
func parseData(ctx *cue.Context, data []byte, format DataFormat, filename string) (cue.Value, error) {
	switch format {
	case FormatJSON:
		return parseJSON(ctx, data, filename)
	case FormatYAML:
		return parseYAML(ctx, data, filename)
	default:
		return cue.Value{}, fmt.Errorf("unsupported format: %d", format)
	}
}

// parseJSON parses JSON data into a CUE value
func parseJSON(ctx *cue.Context, data []byte, filename string) (cue.Value, error) {
	expr, err := json.Extract(filename, data)
	if err != nil {
		return cue.Value{}, fmt.Errorf("parsing JSON: %w", err)
	}
	return ctx.BuildExpr(expr), nil
}

// parseYAML parses YAML data into a CUE value
func parseYAML(ctx *cue.Context, data []byte, filename string) (cue.Value, error) {
	file, err := yaml.Extract(filename, data)
	if err != nil {
		return cue.Value{}, fmt.Errorf("parsing YAML: %w", err)
	}
	return ctx.BuildFile(file), nil
}
