package stub

import (
	"context"
	"fmt"

	"github.com/wolftotem4/golava-new/internal/forge"
)

func ForgeRouteHomeRegister(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
	switch args.DBType.Name {
	case "sqlx":
		worker := CopyGoFile("route.home.register.sqlx.stub", "routes/home/register.go")
		return worker(ctx, args)
	case "gorm":
		worker := CopyGoFile("route.home.register.gorm.stub", "routes/home/register.go")
		return worker(ctx, args)
	}

	return "", nil, fmt.Errorf("db type %s not supported", args.DBType.Name)
}
