package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"sort"
)

func usage(status int) {
	fmt.Printf(`Usage:
	%s [PKG]
where PKG is the name of a Go package (e.g., github.com/cespare/deplist). If no
package name is given, the current directory is used.
`, os.Args[0])
	os.Exit(status)
}

type ctx struct {
	soFar map[string]bool
	cwd   string
}

func (c *ctx) find(name string, testImports bool) error {
	if name == "C" {
		return nil
	}
	pkg, err := build.Import(name, c.cwd, 0)
	if err != nil {
		return err
	}
	if pkg.Goroot {
		return nil
	}

	c.soFar[pkg.ImportPath] = true
	imports := pkg.Imports
	if testImports {
		imports = append(imports, pkg.TestImports...)
	}
	for _, imp := range imports {
		if !c.soFar[imp] {
			if err := c.find(imp, testImports); err != nil {
				return err
			}
		}
	}
	return nil
}

func FindDeps(name string, testImports bool) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	c := &ctx{
		soFar: make(map[string]bool),
		cwd:   cwd,
	}
	if err := c.find(name, testImports); err != nil {
		return nil, err
	}
	var deps []string
	for p := range c.soFar {
		if p != name {
			deps = append(deps, p)
		}
	}
	sort.Strings(deps)
	return deps, nil
}

func main() {
	testImports := flag.Bool("t", false, "Include test dependencies")
	flag.Parse()

	pkg := "."
	switch flag.NArg() {
	case 1:
		pkg = flag.Arg(0)
	default:
		usage(1)
	}

	deps, err := FindDeps(pkg, *testImports)
	if err != nil {
		log.Fatal(err)
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
