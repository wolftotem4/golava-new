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
	fmt.Println("Running go generate ./ent...")

	cmd := exec.Command("go", "generate", "./ent")
	cmd.Dir = dir
	return errors.WithStack(cmd.Run())
}
