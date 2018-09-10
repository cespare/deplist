package main

import (
	"testing"
)

var testCases = []struct {
	name   string
	dir    string
	o      opts
	output []string
}{
	{"c", "/", 0, nil},
	{"b", "/", 0, []string{"c"}},
	{"a", "/", 0, []string{"b", "c"}},
	{".", "testdata/src/a", 0, []string{"b", "c"}},
	{"a", "/", optTestImports, []string{"b", "c", "d"}},
	{"a", "/", optStd, []string{"b", "c", "unsafe"}},
	{"a", "/", optTestImports | optStd, []string{"b", "c", "d", "unsafe"}},
	{"e", "/", 0, []string{"e/vendor/v0", "e/vendor/v0/vendor/a"}},
}

func TestFindDeps(t *testing.T) {
	// FIXME: In the GOPATH world, the test worked by creating a tiny little
	// GOPATH of files and running go/build functions in that restricted
	// gopath context. What's the best way to write these tests in a
	// go/packages world?
	t.Fatal("How to fix?")

	//cwd, err := os.Getwd()
	//if err != nil {
	//        t.Fatal(err)
	//}
	//for _, tt := range testCases {
	//        deps, err := findDeps(tt.name, tt.dir, filepath.Join(cwd, "testdata"), tt.o)
	//        if err != nil {
	//                t.Fatal(err)
	//        }
	//        if !reflect.DeepEqual(deps, tt.output) {
	//                t.Errorf("got %v; want %v\n", deps, tt.output)
	//        }
	//}
}
