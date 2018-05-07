package deplist_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"github.com/elgohr/deplist/deplist"
)

var testCases = []struct {
	name   string
	dir    string
	o      deplist.Opts
	output []string
}{
	{"c", "/", 0, nil},
	{"b", "/", 0, []string{"c"}},
	{"a", "/", 0, []string{"b", "c"}},
	{".", "testdata/src/a", 0, []string{"b", "c"}},
	{"a", "/", deplist.OptTestImports, []string{"b", "c", "d"}},
	{"a", "/", deplist.OptStd, []string{"b", "c", "unsafe"}},
	{"a", "/", deplist.OptTestImports | deplist.OptStd, []string{"b", "c", "d", "unsafe"}},
	{"e", "/", 0, []string{"e/vendor/v0", "e/vendor/v0/vendor/a"}},
}

func TestFindDeps(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range testCases {
		deps, err := deplist.FindDeps(tt.name, tt.dir, filepath.Join(cwd, "testdata"), tt.o)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(deps, tt.output) {
			t.Errorf("got %v; want %v\n", deps, tt.output)
		}
	}
}
