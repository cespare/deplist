package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	tests   = true
	notests = false
)

var testCases = []struct {
	name   string
	dir    string
	tests  bool
	output []string
}{
	{"c", "/", notests, nil},
	{"b", "/", notests, []string{"c"}},
	{"a", "/", notests, []string{"b", "c"}},
	{".", "testdata/src/a", notests, []string{"b", "c"}},
	{"a", "/", tests, []string{"b", "c", "d"}},
}

func TestFindDeps(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range testCases {
		deps, err := FindDeps(tt.name, tt.dir, filepath.Join(cwd, "testdata"), tt.tests)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(deps, tt.output) {
			t.Errorf("got %v; want %v\n", deps, tt.output)
		}
	}
}
