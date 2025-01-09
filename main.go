package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wolftotem4/golava-new/internal/cli/cloneproj"
	"github.com/wolftotem4/golava-new/internal/cli/dotenv"
	"github.com/wolftotem4/golava-new/internal/cli/gomod"
	"github.com/wolftotem4/golava-new/internal/cli/setup"
	"github.com/wolftotem4/golava-new/internal/db"
	"golang.org/x/mod/module"
)

var remoteProj = cloneproj.CloneProject{
	Remote:  "https://github.com/wolftotem4/golava/archive/refs/tags/%s.zip",
	Version: "v0.1.10",
}

var migrations = []string{
	"1732165783_users.down.sql",
	"1732165783_users.up.sql",
	"1732170890_session.down.sql",
	"1732170890_session.up.sql",
}

//go:embed setup.yaml.example
var setupExample embed.FS

var generate = flag.Bool("generate", false, "Generate setup.yaml file")

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: golava-new <project> <moudle-path>\n")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *generate {
		generateSetupYaml()
		return
	}

	if len(flag.Args()) < 2 {
		usage()
		return
	}

	ctx := context.Background()

	err := readSetupConfig(func(config setup.SetupConfig) error {
		run(ctx, config)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func run(ctx context.Context, config setup.SetupConfig) {
	// dir := "../golava"
	// dbType := db.Data["sqlx"]
	// modulePath := "example.com/golava"

	dir := getProjectDir()
	dbType := getDBType(config)
	modulePath := getMoudlePath()

	createProject(dir)
	setupDotEnvFile(dir, config)
	forgeGoFiles(ctx, dir, dbType, false)
	replaceModulePath(dir, modulePath)

	switch dbType.Name {
	case "ent":
		err := gomod.RunGoGenerateEnt(dir)
		if err != nil {
			fmt.Printf("go generate ./ent failed: %s\n", err.Error())
			return
		}
	}

	err := gomod.RunGoModTidy(dir)
	if err != nil {
		fmt.Printf("go mod tidy failed: %s\n", err.Error())
		return
	}

	fmt.Println("Project created successfully")
}

func readSetupConfig(reader func(config setup.SetupConfig) error) error {
	var setupFile io.ReadCloser
	if _, err := os.Stat("setup.yaml"); errors.Is(err, os.ErrNotExist) {
		example, err := setupExample.Open("setup.yaml.example")
		if err != nil {
			return fmt.Errorf("Open setup.yaml.example failed: %s", err.Error())
		}
		defer example.Close()

		fmt.Println("setup.yaml not found, creating one")
		dest, err := os.Create("setup.yaml")
		if err != nil {
			return fmt.Errorf("Create setup.yaml failed: %s", err.Error())
		}
		defer dest.Close()

		setupFile = io.NopCloser(io.TeeReader(example, dest))
	} else {
		setupFile, err = os.Open("setup.yaml")
		if err != nil {
			return fmt.Errorf("Open setup.yaml failed: %s", err.Error())
		}
		defer setupFile.Close()
	}

	config, err := setup.LoadSetupConfig(setupFile)
	if err != nil {
		return fmt.Errorf("Load setup config failed: %s", err.Error())
	}

	return reader(config)
}

func generateSetupYaml() {
	example, err := setupExample.Open("setup.yaml.example")
	if err != nil {
		fmt.Printf("Open setup.yaml.example failed: %s\n", err.Error())
		os.Exit(1)
	}
	defer example.Close()

	dest, err := os.Create("setup.yaml")
	if err != nil {
		fmt.Printf("Create setup.yaml failed: %s\n", err.Error())
		os.Exit(1)
	}
	defer dest.Close()

	_, err = io.Copy(dest, example)
	if err != nil {
		fmt.Printf("Copy setup.yaml.example failed: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("setup.yaml file generated")
	os.Exit(0)
}

func setupDotEnvFile(dir string, config setup.SetupConfig) {
	fmt.Println("Setting up .env file...")
	err := dotenv.SetupDotEnvFile(dir, config)
	if err != nil {
		fmt.Printf("Setup .env file failed: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(".env file setup")
}

func forgeGoFiles(ctx context.Context, dir string, dbType db.DBType, AskOverwrite bool) {
	fmt.Println("Forging go files...")
	files, jobs := prepareForges(ctx, dir, dbType)
	if AskOverwrite {
		promptOverwriteConfirmation(files)
	}
	createFolders(files)
	doForges(ctx, jobs)
	fmt.Println("Go files forged")
}

func createFolders(files []string) {
	var folders = make(map[string]struct{})
	for _, file := range files {
		folders[filepath.Dir(file)] = struct{}{}
	}
	for path := range folders {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func getDBType(config setup.SetupConfig) db.DBType {
	dbType, ok := db.Data[config.DB.Type]
	if !ok {
		fmt.Printf("DB package %s not supported", config.DB.Type)
		os.Exit(1)
	}
	return dbType
}

func getProjectDir() string {
	dir := strings.TrimSpace(flag.Arg(0))
	if dir == "" {
		usage()
		os.Exit(1)
	}
	return dir
}

func getMoudlePath() string {
	modulePath := strings.TrimSpace(flag.Arg(1))
	if modulePath == "" {
		usage()
		os.Exit(1)
	}

	err := module.CheckPath(modulePath)
	if err != nil {
		fmt.Printf("Invalid module path: %s\n", err.Error())
		os.Exit(1)
	}

	return modulePath
}

func createProject(dir string) {
	fmt.Println("Creating project...")
	err := remoteProj.CreateProject(dir)
	if err != nil {
		fmt.Printf("Create project failed: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Project created")
}

func replaceModulePath(dir, modulePath string) {
	fmt.Println("Replacing module path...")
	err := gomod.ReplaceModulePath(dir, modulePath)
	if err != nil {
		fmt.Printf("Replace module path failed: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Module path replaced")
}
