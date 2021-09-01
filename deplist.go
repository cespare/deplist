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

func main() {
	log.SetFlags(0)
	testDeps := flag.Bool("t", false, "Include test dependencies")
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

	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedDeps,
		Tests: *testDeps,
	}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		log.Fatalln("Error loading packages:", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	depSet := make(map[string]struct{})
	pre := func(pkg *packages.Package) bool {
		depSet[pkg.PkgPath] = struct{}{}
		return true
	}
	packages.Visit(pkgs, pre, nil)
	if !*std {
		cfg.Tests = false
		stdPkgs, err := packages.Load(cfg, "std")
		if err != nil {
			log.Fatalln("Error discovering packages:", err)
		}
		if packages.PrintErrors(pkgs) > 0 {
			os.Exit(1)
		}

		stdSet := make(map[string]struct{})
		pre := func(pkg *packages.Package) bool {
			stdSet[pkg.PkgPath] = struct{}{}
			return true
		}
		packages.Visit(stdPkgs, pre, nil)
		for pkg := range stdSet {
			delete(depSet, pkg)
		}
	}
	for _, pkg := range pkgs {
		delete(depSet, pkg.PkgPath)
	}
	var deps []string
	for dep := range depSet {
		deps = append(deps, dep)
	}
	sort.Strings(deps)
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
