package cuebridge

import (
	"fmt"
	"io"
	"os"
)

// readInput reads data from the specified input source
func readInput(input ValidationInput) ([]byte, error) {
	switch input.SourceType {
	case SourceFile:
		return readFromFile(input.FilePath)
	case SourceReader:
		return readFromReader(input.Reader)
	case SourceBytes:
		return readFromBytes(input.Data)
	default:
		return nil, fmt.Errorf("unknown source type: %d", input.SourceType)
	}
}

// readFromFile reads data from a file
func readFromFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}
	return data, nil
}

// readFromReader reads data from an io.Reader
func readFromReader(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading from reader: %w", err)
	}
	return data, nil
}

// readFromBytes returns the data directly
func readFromBytes(data []byte) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	return data, nil
}
