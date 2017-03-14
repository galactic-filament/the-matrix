package Client

import (
	"encoding/json"
	"strings"
)

// TestOutputLine - a line in the test output
type TestOutputLine struct {
	Message  string `json:"message"`
	Line     string `json:"line"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

// TestOutput - a list of test output lines
type TestOutput struct {
	Results []TestOutputLine `json:"output"`
}

// ParseClientLogs - parses container logs
func ParseClientLogs(containerLogs string) (*TestOutput, error) {
	results := &TestOutput{}
	err := json.NewDecoder(strings.NewReader(containerLogs)).Decode(&results)
	if err != nil {
		return &TestOutput{}, err
	}

	return results, nil
}
