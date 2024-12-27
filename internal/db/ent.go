package db

import (
	"fmt"

	"github.com/wolftotem4/golava-new/internal/app"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

var DBTypeEnt = DBType{
	Name:                "ent",
	Handler:             "*sql.DB",
	Package:             pkg.PackageImport{Path: "database/sql"},
	UserProviderPackage: pkg.PackageImport{Path: "github.com/wolftotem4/golava/internal/entauth"},
	UserProvider: `&entauth.EntUserProvider{
			Hasher: app.Hashing,
			Ent:    app.Ent,
		}`,
	MapDBDriver: MapDBDriver{
		"sqlite":   {Package: pkg.PackageImport{Alias: "_", Path: "modernc.org/sqlite"}, Code: sqlDBconn},
		"mysql":    {Package: pkg.PackageImport{Alias: "_", Path: "github.com/go-sql-driver/mysql"}, Code: sqlDBconn},
		"postgres": {Package: pkg.PackageImport{Alias: "_", Path: "github.com/lib/pq"}, Code: sqlDBconn},
	},
	MapDBSessionHandler: MapDBSessionHandler{
		"sqlite": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/sqlite"},
			Code:    fmt.Sprintf(sqlDBSessionHandler, "&sess.SqliteSessionHandler"),
		},
		"mysql": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/mysql"},
			Code:    fmt.Sprintf(sqlDBSessionHandler, "&sess.MySQLSessionHandler"),
		},
		"postgres": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/postgres"},
			Code:    fmt.Sprintf(sqlDBSessionHandler, "&sess.PostgresSessionHandler"),
		},
	},
	AppExts: []app.AppExt{
		{
			Packages: pkg.PackageImports{{Path: "github.com/wolftotem4/golava/ent"}},
			InitPkgs: pkg.PackageImports{{Alias: "entsql", Path: "entgo.io/ent/dialect/sql"}},
			Declare:  "*ent.Client",
			Name:     "Ent",
			Init:     `ent.NewClient(ent.Driver(entsql.OpenDB(getEntDBDriver(), db)))`,
		},
	},
}
