package app

import "github.com/wolftotem4/golava-new/internal/pkg"

type AppExt struct {
	Packages pkg.PackageImports
	Declare  string
	Name     string
	InitPkgs pkg.PackageImports
	Init     string
}

type AppExts []AppExt
