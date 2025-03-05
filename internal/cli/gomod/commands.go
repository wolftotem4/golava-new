package gomod

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func RunGoModTidy(dir string) error {
	fmt.Println("Running go mod tidy...")

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	return errors.WithStack(cmd.Run())
}

func RunGoGenerateEnt(dir string) error {
	fmt.Println("Running go generate ./database/ent...")

	cmd := exec.Command("go", "generate", "./database/ent")
	cmd.Dir = dir
	return errors.WithStack(cmd.Run())
}

func RunGoGet(dir string, pkg string) error {
	fmt.Printf("Running go get %s...\n", pkg)

	cmd := exec.Command("go", "get", pkg)
	cmd.Dir = dir
	return errors.WithStack(cmd.Run())
}
