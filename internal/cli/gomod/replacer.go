package gomod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const MODULE_PATH = "github.com/wolftotem4/golava"

func ReplaceModulePath(dir string, modulePath string) error {
	err := replaceContent(filepath.Join(dir, "go.mod"), fmt.Sprintf("module %s", MODULE_PATH), fmt.Sprintf("module %s", modulePath))
	if err != nil {
		return err
	}

	findme := fmt.Sprintf("%s/", MODULE_PATH)
	replacement := fmt.Sprintf("%s/", modulePath)
	err = filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".go" {
			return nil
		}

		err = replaceContent(path, findme, replacement)
		if err != nil {
			return err
		}

		return nil
	})
	return errors.WithStack(err)
}

func replaceContent(path, old, new string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	plaintext := strings.ReplaceAll(string(content), old, new)

	err = os.WriteFile(path, []byte(plaintext), os.ModePerm)
	return errors.WithStack(err)
}
