package deplist

import (
	"sort"
	"go/build"
)

type Opts uint

const (
	OptTestImports Opts = 1 << iota
	OptStd
)

type context struct {
	soFar map[string]bool
	ctx   build.Context
}

func (c *context) find(name, dir string, o Opts) error {
	if name == "C" {
		return nil
	}
	pkg, err := c.ctx.Import(name, dir, 0)
	if err != nil {
		return err
	}
	if pkg.Goroot && o&OptStd == 0 {
		return nil
	}

	if name != "." {
		c.soFar[pkg.ImportPath] = true
	}
	imports := pkg.Imports
	if o&OptTestImports != 0 {
		imports = append(imports, pkg.TestImports...)
	}
	for _, imp := range imports {
		if !c.soFar[imp] {
			if err := c.find(imp, pkg.Dir, o); err != nil {
				return err
			}
		}
	}
	return nil
}

func FindDeps(name, dir, gopath string, o Opts) ([]string, error) {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	c := &context{
		soFar: make(map[string]bool),
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
