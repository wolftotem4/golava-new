package db

import (
	"fmt"

	"github.com/wolftotem4/golava-new/internal/pkg"
)

const gormDBConn = `db, err := gorm.Open(%s.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		return nil, err
	}`

var gormDBSessionHandler = `native, err := db.DB()
	if err != nil {
		return nil, err
	}
	handler := %s(native, table)`

var DBTypeGORM = DBType{
	Name:                "gorm",
	Handler:             "*gorm.DB",
	Package:             pkg.PackageImport{Path: "gorm.io/gorm"},
	UserProviderPackage: pkg.PackageImport{Alias: "db", Path: "github.com/wolftotem4/golava-db-gorm"},
	UserProvider: `&db.GormUserProvider{
			Hasher:        a.Hashing,
			DB:            a.DB,
			ConstructUser: func() auth.Authenticatable { return &generic.User{} },
		}`,
	MapDBDriver: MapDBDriver{
		"sqlite":    {Package: pkg.PackageImport{Path: "gorm.io/driver/sqlite"}, Code: fmt.Sprintf(gormDBConn, "sqlite")},
		"mysql":     {Package: pkg.PackageImport{Path: "gorm.io/driver/mysql"}, Code: fmt.Sprintf(gormDBConn, "mysql")},
		"postgres":  {Package: pkg.PackageImport{Path: "gorm.io/driver/postgres"}, Code: fmt.Sprintf(gormDBConn, "postgres")},
		"sqlserver": {Package: pkg.PackageImport{Path: "gorm.io/driver/sqlserver"}, Code: fmt.Sprintf(gormDBConn, "sqlserver")},
	},
	MapDBSessionHandler: MapDBSessionHandler{
		"sqlite": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/sqlite"},
			Code:    fmt.Sprintf(gormDBSessionHandler, "sess.NewSqliteSessionHandler"),
		},
		"mysql": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/mysql"},
			Code:    fmt.Sprintf(gormDBSessionHandler, "sess.NewMySQLSessionHandler"),
		},
		"postgres": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/postgres"},
			Code:    fmt.Sprintf(gormDBSessionHandler, "sess.NewPostgresSessionHandler"),
		},
		"sqlserver": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/sqlserver"},
			Code:    fmt.Sprintf(gormDBSessionHandler, "sess.NewSQLServerSessionHandler"),
		},
	},
}
