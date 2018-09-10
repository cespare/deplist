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

where pkg is the name of a Go package (e.g., github.com/cespare/deplist).
If no package name is given, the current directory is used.

Flags:
`, os.Args[0])
	flag.PrintDefaults()
}

type opts uint

const (
	optTestImports opts = 1 << iota
	optStd
)

type context struct {
	soFar map[string]struct{}
	ctx   build.Context
}

func (c *context) find(name, dir string, o opts) error {
	if name == "C" {
		return nil
	}
	pkg, err := c.ctx.Import(name, dir, 0)
	if err != nil {
		return err
	}
	if pkg.Goroot && o&optStd == 0 {
		return nil
	}

	if name != "." {
		c.soFar[pkg.ImportPath] = struct{}{}
	}
	imports := pkg.Imports
	if o&optTestImports != 0 {
		imports = append(imports, pkg.TestImports...)
	}
	for _, imp := range imports {
		if _, ok := c.soFar[imp]; !ok {
			if err := c.find(imp, pkg.Dir, o); err != nil {
				return err
			}
		}
	}
	return nil
}

func findDeps(name, dir, gopath string, o opts) ([]string, error) {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	c := &context{
		soFar: make(map[string]struct{}),
		ctx:   ctx,
	}
	if err := c.find(name, dir, o); err != nil {
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
	std := flag.Bool("std", false, "Include standard library dependencies")
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

	var o opts
	if *testImports {
		o |= optTestImports
	}
	if *std {
		o |= optStd
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't determine working directory:", err)
	}
	deps, err := findDeps(pkg, cwd, "", o)
	if err != nil {
		log.Fatal(err)
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
