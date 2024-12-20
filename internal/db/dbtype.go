package db

import "github.com/wolftotem4/golava-new/internal/pkg"

type DBType struct {
	Name                string
	Handler             string
	Package             pkg.PackageImport
	UserProviderPackage pkg.PackageImport
	UserProvider        string
	MapDBDriver         MapDBDriver
	MapDBSessionHandler MapDBSessionHandler
}

var Data = map[string]DBType{
	DBTypeSQLX.Name: DBTypeSQLX,
	DBTypeGORM.Name: DBTypeGORM,
}