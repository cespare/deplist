package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
)

func usage(status int) {
	fmt.Printf(`Usage:
	%s [PKG]
where PKG is the name of a Go package (e.g., github.com/cespare/deplist). If no
package name is given, the current directory is used.
`, os.Args[0])
	os.Exit(status)
}

func findDeps(soFar map[string]bool, name string, testImports bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if name == "C" {
		return nil
	}
	pkg, err := build.Import(name, cwd, 0)
	if err != nil {
		return err
	}
	if pkg.Goroot {
		return nil
	}

	soFar[pkg.ImportPath] = true
	imports := pkg.Imports
	if testImports {
		imports = append(imports, pkg.TestImports...)
	}
	for _, imp := range imports {
		if !soFar[imp] {
			if err := findDeps(soFar, imp, testImports); err != nil {
				return err
			}
		}
	}
	return nil
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

	deps := make(map[string]bool)
	err := findDeps(deps, pkg, *testImports)
	if err != nil {
		log.Fatalln(err)
	}
	delete(deps, pkg)
	for dep := range deps {
		fmt.Println(dep)
	}
}
