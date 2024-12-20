package db

import "github.com/wolftotem4/golava-new/internal/pkg"

type DBSessionHandler struct {
	Package pkg.PackageImport
	Code    string
}

type MapDBSessionHandler map[string]DBSessionHandler
