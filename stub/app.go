package stub

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

func ForgeAppGo(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = filepath.Join(args.Dir, "internal/app/app.go")

	return file, func(ctx context.Context) error {
		demo, _ := os.Create(file)
		defer demo.Close()

		packages := pkg.PackageImports{
			{Path: "github.com/wolftotem4/golava-core/golava"},
			{Path: "github.com/wolftotem4/golava/internal/logging"},
			args.DBType.Package,
		}

		for _, ext := range args.DBType.AppExts {
			packages = append(packages, ext.Packages...)
		}

		packages.Unique()
		sort.Sort(packages)

		return parseTemplate("app.db.stub", demo, map[string]any{
			"packages": packages.String(),
			"dbType":   args.DBType.Handler,
			"exts":     args.DBType.AppExts,
		})
	}, nil
}
