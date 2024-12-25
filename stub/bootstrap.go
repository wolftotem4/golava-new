package stub

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

func ForgeBootstrap(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = filepath.Join(args.Dir, "internal/bootstrap/bootstrap.go")

	return file, func(ctx context.Context) error {
		code, _ := os.Create(file)
		defer code.Close()

		driverName := args.DBDriver

		packages := pkg.PackageImports{
			{Path: "context"},
			{Path: "log/slog"},
			{Path: "os"},
			{Path: "time"},
			{Path: "github.com/joho/godotenv"},
			{Path: "github.com/wolftotem4/golava-core/cookie"},
			{Path: "github.com/wolftotem4/golava-core/golava"},
			{Path: "github.com/wolftotem4/golava-core/hashing"},
			{Path: "github.com/wolftotem4/golava-core/router"},
			{Path: "github.com/wolftotem4/golava/internal/app"},
			args.DBType.Package,
			args.DBType.MapDBDriver[driverName].Package,
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

		return parseTemplate("bootstrap.db.stub", code, map[string]any{
			"packages":   packages.String(),
			"dbTypeName": args.DBType.Name,
			"dbConn":     args.DBType.MapDBDriver[driverName].Code,
			"exts":       args.DBType.AppExts,
		})
	}, nil
}
