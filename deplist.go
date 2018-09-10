package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/tools/go/packages"
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

func findDeps(name, dir string, o opts) ([]string, error) {
	config := &packages.Config{
		Mode:  packages.LoadImports,
		Dir:   dir,
		Env:   packagesEnv,
		Tests: o&optTestImports != 0,
	}
	var ignore map[string]struct{}
	if o&optStd != 0 {
		var err error
		ignore, err = findStd()
		if err != nil {
			return nil, fmt.Errorf("error loading std packages: %s", err)
		}
	}
	pkgs, err := packages.Load(config, name)
	if err != nil {
		return nil, fmt.Errorf("error loading specified package(s): %s", err)
	}
	topLevel := make(map[string]struct{})
	for _, pkg := range pkgs {
		topLevel[pkg.PkgPath] = struct{}{}
	}
	depSet := make(map[string]struct{})
	fn := func(p *packages.Package) bool {
		pp := p.PkgPath
		if contains(topLevel, pp) && !contains(ignore, pp) {
			depSet[pp] = struct{}{}
		}
		return true
	}
	packages.Visit(pkgs, fn, nil)
	var deps []string
	for p := range depSet {
		deps = append(deps, p)
	}
	sort.Strings(deps)
	return deps, nil
}

func findStd() (map[string]struct{}, error) {
	config := &packages.Config{
		Env: packagesEnv,
	}
	pkgs, err := packages.Load(config, "std")
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{})
	for _, pkg := range pkgs {
		set[pkg.PkgPath] = struct{}{}
	}
	return set, nil
}

func contains(set map[string]struct{}, s string) bool {
	_, ok := set[s]
	return ok
}

var packagesEnv = append(os.Environ(), "GOFLAGS=-mod=readonly")

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
	// FIXME: How to exclude packages from the stdlib?
	if *std {
		o |= optStd
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't determine working directory:", err)
	}
	deps, err := findDeps(pkg, cwd, o)
	if err != nil {
		log.Fatal(err)
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
