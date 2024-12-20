package db

import "github.com/wolftotem4/golava-new/internal/pkg"

type DBDriver struct {
	Package pkg.PackageImport
	Code    string
}

type MapDBDriver map[string]DBDriver
