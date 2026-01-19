package e2e

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
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

	expected := map[string]string{
		fileA.Name(): fmt.Sprintf(" 1 5 24 %s", fileA.Name()),
		fileB.Name(): fmt.Sprintf(" 2 3 13 %s", fileB.Name()),
		fileC.Name(): fmt.Sprintf(" 0 0  0 %s", fileC.Name()),
		"total":      fmt.Sprintf(" 3 8 37 %s", "total"),
	}

	expectedChecks := len(expected)
	checks := 0
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			t.Log("Encountered an empty line")
			t.Fail()
		}

		filename := fields[len(fields)-1]

		lineExpected, ok := expected[filename]
		if !ok {
			t.Errorf("unexpected filename: %s", filename)
			continue
		}

		if line != lineExpected {
			t.Errorf("unexpected output for %s: got %s, expected %s", filename, line, lineExpected)
			t.Fail()
		}

		checks++
	}

	if checks != expectedChecks {
		t.Errorf("unexpected number of checks: got %d, expected %d", checks, expectedChecks)
		t.Fail()
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("failed to scan stdout: %v", err)
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
