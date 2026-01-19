package e2e

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestMultiFile(t *testing.T) {
	fileA, err := CreateFile("one two three four five\n")
	if err != nil {
		t.Fatalf("failed to create fileA: %v", err)
	}
	defer os.Remove(fileA.Name())

	fileB, err := CreateFile("foo bar baz\n\n")
	if err != nil {
		t.Fatalf("failed to create fileB: %v", err)
	}
	defer os.Remove(fileB.Name())

	fileC, err := CreateFile("")
	if err != nil {
		t.Fatalf("failed to create fileC: %v", err)
	}
	defer os.Remove(fileC.Name())

	cmd, err := GetCommand(fileA.Name(), fileB.Name(), fileC.Name())
	if err != nil {
		t.Fatalf("failed to get command: %v", err)
	}

	stdout := &bytes.Buffer{}

	cmd.Stdout = stdout

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run command: %v", err)
	}

	expectedOut := fmt.Sprintf(
		" 1 5 24 %s\n 2 3 13 %s\n 0 0  0 %s\n 3 8 37 %s\n",
		fileA.Name(),
		fileB.Name(),
		fileC.Name(),
		"total")

	res := stdout.String()
	if res != expectedOut {
		t.Errorf("unexpected output: got %s, expected %s", res, expectedOut)
		t.Fail()
	}
}

func CreateFile(content string) (*os.File, error) {
	file, err := os.CreateTemp("", "counter-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := file.WriteString(content); err != nil {
		return nil, fmt.Errorf("failed to write content to file: %w", err)
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return file, nil
}
