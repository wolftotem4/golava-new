package stub

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/wolftotem4/golava-new/internal/forge"
)

func CopyFile(src string, dest string) forge.ForgeWorker {
	return func(ctx context.Context, args forge.ForgeWorkArgs) (gofile string, forge func(ctx context.Context) error, err error) {
		var file = filepath.Join(args.Dir, dest)

		return file, func(ctx context.Context) error {
			destFh, err := os.Create(file)
			if err != nil {
				return err
			}
			defer destFh.Close()

			srcFh, err := stub.Open(path.Join("files", src))
			if err != nil {
				return err
			}
			defer srcFh.Close()

			_, err = io.Copy(destFh, srcFh)
			return err
		}, nil
	}
}
