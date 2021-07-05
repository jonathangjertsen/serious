package cmd

import (
	"testing"
)

func Test_UsageErrorWithoutAllowed(t *testing.T) {
	err := UsageError{
		Opt:     "foo",
		Value:   "some value",
		Allowed: []string{},
	}
	have := err.Error()
	want := `foo can not be 'some value'`
	if have != want {
		t.Fatalf("have:\n%s\n\nwant:\n%s", have, want)
	}
}

func Test_UsageErrorWithAllowed(t *testing.T) {
	err := UsageError{
		Opt:     "foo",
		Value:   "some value",
		Allowed: []string{"bar", "baz"},
	}
	have := err.Error()
	want := `foo can not be 'some value' (allowed: bar, baz)`
	if have != want {
		t.Fatalf("have:\n%s\n\nwant:\n%s", have, want)
	}
}
