package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal(t *testing.T, expected, actual any, msg ...string) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return
	}

	msgStr := strings.Join(msg, ": ")

	t.Logf("\n%s\nExpected: %v Got: %v", msgStr, expected, actual)
	t.Fail()
}
