package stub

import (
	"context"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/internal/pkg"
)

func ForgeMiddlewareAuth(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	var file = path.Join(args.Dir, "internal/middlewares/auth.go")

	return file, func(ctx context.Context) error {
		code, _ := os.Create(file)
		defer code.Close()

		packages := pkg.PackageImports{
			{Path: "github.com/gin-gonic/gin"},
			{Path: "github.com/wolftotem4/golava-core/auth/generic"},
			{Path: "github.com/wolftotem4/golava-core/instance"},
			{Path: "github.com/wolftotem4/golava/internal/app"},
			args.DBType.UserProviderPackage,
		}

		if hasAuth(args.DBType.UserProvider) {
			packages = append(packages, pkg.PackageImport{Path: "github.com/wolftotem4/golava-core/auth"})
		}

		sort.Sort(packages)

		return parseTemplate("middleware.auth.stub", code, map[string]any{
			"packages":     packages.String(),
			"userProvider": args.DBType.UserProvider,
		})
	}, nil
}

func hasAuth(content string) bool {
	return regexp.MustCompile(`\bauth\.`).MatchString(content)
}
