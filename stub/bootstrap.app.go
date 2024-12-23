package stub

import (
	"context"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

func ForgeBootstrapApp(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = path.Join(args.Dir, "internal/bootstrap/app.go")

	return file, func(ctx context.Context) error {
		code, _ := os.Create(file)
		defer code.Close()

		driverName := args.DBDriver

		packages := pkg.PackageImports{
			{Path: "github.com/joho/godotenv"},
			{Path: "github.com/wolftotem4/golava-core/cookie"},
			{Path: "github.com/wolftotem4/golava-core/encryption"},
			{Path: "github.com/wolftotem4/golava-core/golava"},
			{Path: "github.com/wolftotem4/golava-core/hashing"},
			{Path: "github.com/wolftotem4/golava-core/router"},
			{Path: "github.com/wolftotem4/golava-core/session"},
			{Path: "github.com/wolftotem4/golava/internal/app"},
			args.DBType.Package,
			args.DBType.MapDBDriver[driverName].Package,
			args.DBType.MapDBSessionHandler[driverName].Package,
		}

		for _, ext := range args.DBType.AppExts {
			packages = append(packages, ext.Packages...)
			packages = append(packages, ext.InitPkgs...)
		}

		packages.Unique()
		sort.Sort(packages)

		if driverName == "sqlite" {
			for index, ext := range args.DBType.AppExts {
				if ext.Name == "Ent" {
					args.DBType.AppExts[index].Init = strings.Replace(ext.Init, `os.Getenv("DB_DRIVER")`, `"sqlite3"`, 1)
				}
			}
		}

		return parseTemplate("bootstrap.app.db.stub", code, map[string]any{
			"packages":    packages.String(),
			"dbType":      args.DBType.Handler,
			"dbConn":      args.DBType.MapDBDriver[driverName].Code,
			"sessHandler": args.DBType.MapDBSessionHandler[driverName].Code,
			"exts":        args.DBType.AppExts,
		})
	}, nil
}
