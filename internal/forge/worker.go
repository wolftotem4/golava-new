package forge

import (
	"context"

	"github.com/wolftotem4/golava-new/internal/db"
)

type ForgeWorkArgs struct {
	Dir      string
	DBType   db.DBType
	DBDriver string
}

type ForgeWorker func(ctx context.Context, args ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error)

type ForgeWorkers []ForgeWorker

func (f ForgeWorkers) Ready(ctx context.Context, args ForgeWorkArgs) (gofiles []string, forges []func(ctx context.Context) error, err error) {
	for _, worker := range f {
		gofile, forge, err := worker(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		gofiles = append(gofiles, gofile)
		forges = append(forges, forge)
	}
	return gofiles, forges, nil
}
