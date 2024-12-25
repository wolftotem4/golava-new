package pkg

import (
	"fmt"
	"strings"
)

type PackageImport struct {
	Alias string
	Path  string
}

func (p PackageImport) String() string {
	if p.Alias == "" {
		return fmt.Sprintf("%q", p.Path)
	}
	return fmt.Sprintf("%s %q", p.Alias, p.Path)
}

func (p PackageImport) IsNative() bool {
	return IsNativePackage(p.Path)
}

type PackageImports []PackageImport

func (p PackageImports) Len() int { return len(p) }
func (p PackageImports) Less(i, j int) bool {
	leftN := p[i].IsNative()
	rightN := p[j].IsNative()
	if leftN && !rightN {
		return true
	} else if !leftN && rightN {
		return false
	}

	if p[i].Path == p[j].Path {
		return p[i].Alias < p[j].Alias
	}
	return p[i].Path < p[j].Path
}
func (p PackageImports) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PackageImports) String() string {
	var (
		buf                    strings.Builder
		paddingAfterNativePkgs = false
	)
	for _, pkg := range p {
		if pkg.IsNative() {
			paddingAfterNativePkgs = true
		} else if paddingAfterNativePkgs {
			buf.WriteString("\n")
			paddingAfterNativePkgs = false
		}

		buf.WriteString(fmt.Sprintf("\t%s\n", pkg.String()))
	}
	return buf.String()[:buf.Len()-1]
}

func (p *PackageImports) Add(pkg PackageImport) {
	*p = append(*p, pkg)
}

func (p *PackageImports) Unique() {
	m := make(map[PackageImport]struct{})
	for _, pkg := range *p {
		m[pkg] = struct{}{}
	}

	var unique PackageImports
	for pkg := range m {
		unique = append(unique, pkg)
	}

	*p = unique
}
