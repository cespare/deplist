package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"sort"
)

func usage() {
	fmt.Printf(`Usage:
    %s [flags] [pkg]
where pkg is the name of a Go package (e.g., github.com/cespare/deplist). If no
package name is given, the current directory is used.

Flags:
`, os.Args[0])
	flag.PrintDefaults()
}

type context struct {
	soFar map[string]bool
	ctx   build.Context
	dir   string
}

func (c *context) find(name string, testImports bool) error {
	if name == "C" {
		return nil
	}
	pkg, err := c.ctx.Import(name, c.dir, 0)
	if err != nil {
		return err
	}
	if pkg.Goroot {
		return nil
	}

	if name != "." {
		c.soFar[pkg.ImportPath] = true
	}
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

func FindDeps(name, dir, gopath string, testImports bool) ([]string, error) {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	c := &context{
		soFar: make(map[string]bool),
		ctx:   ctx,
		dir:   dir,
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
	flag.Usage = usage
	flag.Parse()

	pkg := "."
	switch flag.NArg() {
	case 1:
		pkg = flag.Arg(0)
	case 0:
	default:
		usage()
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't determine working directory:", err)
	}
	deps, err := FindDeps(pkg, cwd, "", *testImports)
	if err != nil {
		log.Fatal(err)
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
