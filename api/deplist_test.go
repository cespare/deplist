package api_test

import (
	"github.com/elgohr/deplist/api"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testCases = []struct {
	name   string
	dir    string
	o      api.Opts
	output []string
}{
	{"c", "/", 0, nil},
	{"b", "/", 0, []string{"c"}},
	{"a", "/", 0, []string{"b", "c"}},
	{".", "testdata/src/a", 0, []string{"b", "c"}},
	{"a", "/", api.OptTestImports, []string{"b", "c", "d"}},
	{"a", "/", api.OptStd, []string{"b", "c", "unsafe"}},
	{"a", "/", api.OptTestImports | api.OptStd, []string{"b", "c", "d", "unsafe"}},
	{"e", "/", 0, []string{"e/vendor/v0", "e/vendor/v0/vendor/a"}},
}

func TestFindDeps(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range testCases {
		deps, err := api.FindDeps(tt.name, tt.dir, filepath.Join(cwd, "testdata"), tt.o)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(deps, tt.output) {
			t.Errorf("got %v; want %v\n", deps, tt.output)
		}
	}
}
