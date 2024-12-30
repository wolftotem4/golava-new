package stub

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

func ForgeBootstrapApp(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = filepath.Join(args.Dir, "internal/bootstrap/app.go")

	return file, func(ctx context.Context) error {
		code, _ := os.Create(file)
		defer code.Close()

		driverName := args.DBDriver

		packages := pkg.PackageImports{
			{Path: "context"},
			{Path: "os"},
			{Path: "github.com/wolftotem4/golava-core/cookie"},
			{Path: "github.com/wolftotem4/golava-core/golava"},
			{Path: "github.com/wolftotem4/golava-core/hashing"},
			{Path: "github.com/wolftotem4/golava-core/routing"},
			{Path: "github.com/wolftotem4/golava/internal/app"},
			{Path: "github.com/wolftotem4/golava/internal/env"},
			args.DBType.Package,
			args.DBType.MapDBDriver[driverName].Package,
		}

		for _, ext := range args.DBType.AppExts {
			packages = append(packages, ext.Packages...)
			packages = append(packages, ext.InitPkgs...)
		}

		switch args.DBType.Name {
		case "sqlx", "ent":
			packages = append(packages, pkg.PackageImport{Path: "time"})
		}

		packages.Unique()
		sort.Sort(packages)

		return parseTemplate("bootstrap.app.stub", code, map[string]any{
			"packages":   packages.String(),
			"dbTypeName": args.DBType.Name,
			"dbConn":     args.DBType.MapDBDriver[driverName].Code,
			"exts":       args.DBType.AppExts,
		})
	}, nil
}
