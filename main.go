package main

import (
	"flag"
	"fmt"
	"github.com/elgohr/deplist/api"
	"log"
	"os"
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

	var o api.Opts
	if *testImports {
		o |= api.OptTestImports
	}
	if *std {
		o |= api.OptStd
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Couldn't determine working directory:", err)
	}
	deps, err := api.FindDeps(pkg, cwd, "", o)
	if err != nil {
		log.Fatal(err)
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
}
