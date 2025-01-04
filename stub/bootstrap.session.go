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

func ForgeBootstrapSession(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = filepath.Join(args.Dir, "internal/bootstrap/session.go")

	return file, func(ctx context.Context) error {
		code, _ := os.Create(file)
		defer code.Close()

		driverName := args.DBDriver

		packages := pkg.PackageImports{
			{Path: "os"},
			{Path: "time"},
			{Path: "github.com/wolftotem4/golava-core/session"},
			{Path: "github.com/wolftotem4/golava/internal/env"},
			args.DBType.Package,
			args.DBType.MapDBSessionHandler[driverName].Package,
		}
		sort.Sort(packages)

		if driverName == "sqlite" {
			for index, ext := range args.DBType.AppExts {
				if ext.Name == "Ent" {
					args.DBType.AppExts[index].Init = strings.Replace(ext.Init, `os.Getenv("DB_DRIVER")`, `"sqlite3"`, 1)
				}
			}
		}

		return parseTemplate("bootstrap.session.stub", code, map[string]any{
			"packages":    packages.String(),
			"dbType":      args.DBType.Handler,
			"sessHandler": args.DBType.MapDBSessionHandler[driverName].Code,
		})
	}, nil
}
