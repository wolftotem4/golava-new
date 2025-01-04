package dotenv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/wolftotem4/golava-new/internal/cli"
)

const DotEnvExample = ".env.example"

func CreateDotEnvFileAndLoad(dir string) error {
	file := filepath.Join(dir, DotEnvFile)

	if _, err := os.Stat(file); err == nil {
		overwrite := ConfirmDotEnvOverwrite()
		if !overwrite {
			return errors.WithStack(cli.ErrOverwriteRejected)
		}

		if err := CopyEnvFile(dir, true); err != nil {
			return err
		}
	} else {
		if err := CopyEnvFile(dir, false); err != nil {
			return err
		}
	}

	if err := godotenv.Load(file); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func CopyEnvFile(dir string, overwrite bool) error {
	src := filepath.Join(dir, DotEnvExample)
	dst := filepath.Join(dir, DotEnvFile)

	if _, err := os.Stat(dst); err == nil && !overwrite {
		return errors.WithStack(os.ErrExist)
	} else if err != nil && !os.IsNotExist(err) {
		return errors.WithStack(err)
	}

	in, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return errors.WithStack(err)
}

func ConfirmDotEnvOverwrite() bool {
	fmt.Print("Do you want to overwrite the existing .env file? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	return answer == "y"
}

func ReloadDotEnv(dir string) error {
	return errors.WithStack(godotenv.Overload(filepath.Join(dir, DotEnvFile)))
}
