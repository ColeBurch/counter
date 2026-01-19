package e2e

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ColeBurch/counter/test/assert"
)

func TestFlags(t *testing.T) {
	file, err := CreateFile("one two three\n four five \nsix seven\n")
	if err != nil {
		t.Fatal("Failed to create file")
	}
	defer os.Remove(file.Name())

	t.Run("Line flag", func(t *testing.T) {
		cmd, err := GetCommand("-l", file.Name())
		if err != nil {
			t.Fatal("Failed to create command")
			t.Fail()
		}

		stdout := &bytes.Buffer{}
		cmd.Stdout = stdout
		err = cmd.Run()
		if err != nil {
			t.Fatal("Failed to run command")
		}

		expected := fmt.Sprintf(" 3 %s\n 3 total\n", file.Name())
		assert.Equal(t, expected, stdout.String())
	})

	t.Run("Word flag", func(t *testing.T) {
		cmd, err := GetCommand("-w", file.Name())
		if err != nil {
			t.Fatal("Failed to create command")
			t.Fail()
		}

		stdout := &bytes.Buffer{}
		cmd.Stdout = stdout
		err = cmd.Run()
		if err != nil {
			t.Fatal("Failed to run command")
		}

		expected := fmt.Sprintf(" 7 %s\n 7 total\n", file.Name())
		assert.Equal(t, expected, stdout.String())
	})

	t.Run("Character flag", func(t *testing.T) {
		cmd, err := GetCommand("-c", file.Name())
		if err != nil {
			t.Fatal("Failed to create command")
			t.Fail()
		}

		stdout := &bytes.Buffer{}
		cmd.Stdout = stdout
		err = cmd.Run()
		if err != nil {
			t.Fatal("Failed to run command")
		}

		expected := fmt.Sprintf(" 36 %s\n 36 total\n", file.Name())
		assert.Equal(t, expected, stdout.String())
	})
}
