package stub

import (
	"context"
	"fmt"

	"github.com/wolftotem4/golava-new/internal/forge"
)

func ForgeRouteHomeRegister(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	const target = "internal/routes/home/register.go"

	switch args.DBType.Name {
	case "sqlx":
		worker := CopyFile("route.home.register.sqlx.stub", target)
		return worker(ctx, args)
	case "gorm":
		worker := CopyFile("route.home.register.gorm.stub", target)
		return worker(ctx, args)
	case "ent":
		worker := CopyFile("route.home.register.ent.stub", target)
		return worker(ctx, args)
	}

	return "", nil, fmt.Errorf("db type %s not supported", args.DBType.Name)
}
