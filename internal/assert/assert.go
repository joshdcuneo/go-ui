package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("actual %v; expected %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expected string) {
	t.Helper()
	if !strings.Contains(actual, expected) {
		t.Errorf("actual %v; expected to contain %v", actual, expected)
	}
}

func StringNotContains(t *testing.T, actual, expected string) {
	t.Helper()
	if strings.Contains(actual, expected) {
		t.Errorf("actual %v; expected not to contain %v", actual, expected)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("actual %v; expected nil", actual)
	}
}
