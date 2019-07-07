package textwrap_test

import (
	"strings"
	"testing"

	"github.com/skillian/textwrap"
)

const (
	unwrapped = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

	wrapped = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.`
)

func TestString(t *testing.T) {
	t.Parallel()

	result := textwrap.String(unwrapped, 80)

	if result != wrapped {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n\n", wrapped, result)
	} else {
		t.Log("OK")
	}
}

func TestSliceLines(t *testing.T) {
	t.Parallel()

	expected := strings.Split(wrapped, "\n")

	result := textwrap.SliceLines(strings.Split(unwrapped, " "), 80, " ")

	if len(result) != len(expected) {
		t.Fatal("Expected", len(expected), "lines but got", len(result), "lines")
	}

	for i, v := range expected {
		if v != result[i] {
			t.Fatalf(
				"Line %d:  Expected:\n\n\t\t%s\n\n\tbut "+
					"found:\n\n\t\t%s",
				i+1, v, result[i])
		}
	}

}
