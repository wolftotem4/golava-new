package db

import (
	"fmt"

	"github.com/wolftotem4/golava-new/internal/pkg"
)

const sqlxDBconn = `db, err := sqlx.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)`

var sqlxDBSessionHandler = `handler := %s(db.DB, table)`

var DBTypeSQLX = DBType{
	Name:                "sqlx",
	Handler:             "*sqlx.DB",
	Package:             pkg.PackageImport{Path: "github.com/jmoiron/sqlx"},
	UserProviderPackage: pkg.PackageImport{Alias: "db", Path: "github.com/wolftotem4/golava-db-sqlx"},
	UserProvider: `&db.SqlxUserProvider{
			Hasher:        a.Hashing,
			Table:         "users",
			DB:            a.DB,
			ConstructUser: func() auth.Authenticatable { return &generic.User{} },
		}`,
	MapDBDriver: MapDBDriver{
		"sqlite":    {Package: pkg.PackageImport{Alias: "_", Path: "modernc.org/sqlite"}, Code: sqlxDBconn},
		"mysql":     {Package: pkg.PackageImport{Alias: "_", Path: "github.com/go-sql-driver/mysql"}, Code: sqlxDBconn},
		"postgres":  {Package: pkg.PackageImport{Alias: "_", Path: "github.com/lib/pq"}, Code: sqlxDBconn},
		"sqlserver": {Package: pkg.PackageImport{Alias: "_", Path: "github.com/microsoft/go-mssqldb"}, Code: sqlxDBconn},
	},
	MapDBSessionHandler: MapDBSessionHandler{
		"sqlite": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/sqlite"},
			Code:    fmt.Sprintf(sqlxDBSessionHandler, "sess.NewSqliteSessionHandler"),
		},
		"mysql": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/mysql"},
			Code:    fmt.Sprintf(sqlxDBSessionHandler, "sess.NewMySQLSessionHandler"),
		},
		"postgres": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/postgres"},
			Code:    fmt.Sprintf(sqlxDBSessionHandler, "sess.NewPostgresSessionHandler"),
		},
		"sqlserver": {
			Package: pkg.PackageImport{Alias: "sess", Path: "github.com/wolftotem4/golava-core/session/sqlserver"},
			Code:    fmt.Sprintf(sqlxDBSessionHandler, "sess.NewSQLServerSessionHandler"),
		},
	},
}
