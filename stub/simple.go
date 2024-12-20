package stub

import (
	"context"
	"os"
	"path"

	"github.com/wolftotem4/golava-new/internal/forge"
)

func CopyGoFile(src string, dest string) forge.ForgeWorker {
	return func(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
		var file = path.Join(args.Dir, dest)

		return file, func(ctx context.Context) error {
			code, _ := os.Create(file)
			defer code.Close()

			return parseTemplate(src, code, nil)
		}, nil
	}
}
