package cloneproj

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/wolftotem4/golava-new/internal/cli"
)

type CloneProject struct {
	Remote  string
	Version string
}

func (cp CloneProject) GetRemote() string {
	return fmt.Sprintf(cp.Remote, cp.Version)
}

func (cp CloneProject) CreateProject(dir string) error {
	zipFile := fmt.Sprintf("golava-%s.zip", cp.Version)

	err := cli.Download(cp.GetRemote(), zipFile)
	if err != nil {
		return err
	}

	err = unzipGithubProject(dir, zipFile)
	if err != nil {
		return err
	}

	err = os.Remove(zipFile)
	return errors.WithStack(err)
}

func unzipGithubProject(dir string, zipFile string) error {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer archive.Close()

	pathPrefix := archive.File[0].Name

	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		src, err := file.Open()
		if err != nil {
			return errors.WithStack(err)
		}
		defer src.Close()

		path := strings.TrimPrefix(file.Name, pathPrefix)

		err = os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}

		dst, err := os.Create(filepath.Join(dir, path))
		if err != nil {
			return errors.WithStack(err)
		}
		defer dst.Close()

		_, err = dst.ReadFrom(src)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
