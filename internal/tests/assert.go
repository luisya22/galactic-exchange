package assert

import (
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func Greater[T constraints.Ordered](t *testing.T, actual, expected T) {
	t.Helper()

	if !(actual > expected) {
		t.Errorf("got: %v; want greater than: %v", actual, expected)
	}
}

func Smaller[T constraints.Ordered](t *testing.T, actual, expected T) {
	t.Helper()

	if !(actual < expected) {
		t.Errorf("got: %v; want greater than: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %v; expected to contain: %v", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func Error(t *testing.T, actual error) {
	t.Helper()

	if actual == nil {
		t.Errorf("got: %v; expected: not nil", actual)
	}
}
