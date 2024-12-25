package main

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava-new/internal/cli"
	"github.com/wolftotem4/golava-new/internal/cli/question"
	"github.com/wolftotem4/golava-new/internal/db"
	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/stub"
)

const GIT_REMOTE = "https://github.com/wolftotem4/golava.git"
const VERSION = "v0.1.5"

var migrations = []string{
	"1732165783_users.down.sql",
	"1732165783_users.up.sql",
	"1732170890_session.down.sql",
	"1732170890_session.up.sql",
}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Println("Usage: golava-new <project>")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// dir := "../golava"
	// dbType := db.Data["sqlx"]

	dir := getProjectDir()

	createProject(dir)
	setupDotEnvFile(dir)
	dbType := askDBType()
	forgeGoFiles(ctx, dir, dbType)

	switch dbType.Name {
	case "ent":
		runGoGenerateEnt(dir)
	}

	runGoModTidy(dir)

	fmt.Println("Project created successfully")
}

func setupDotEnvFile(dir string) {
	createDotEnvFileAndLoad(dir)
	configureDotEnv(dir)
	reloadDotEnv(dir)
}

func forgeGoFiles(ctx context.Context, dir string, dbType db.DBType) {
	files, jobs := prepareForges(ctx, dir, dbType)
	askOverwrite(files)
	createFolders(files)
	doForges(ctx, jobs)
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

func doForges(ctx context.Context, jobs []func(ctx context.Context) error) {
	for _, job := range jobs {
		if err := job(ctx); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}
	}
}

func askOverwrite(files []string) {
	if !question.AskOverwrite(files) {
		fmt.Println("Operation aborted")
		os.Exit(1)
	}
}

func prepareForges(ctx context.Context, dir string, dbType db.DBType) ([]string, []func(ctx context.Context) error) {
	args := forge.ForgeWorkArgs{
		Dir:      dir,
		DBType:   dbType,
		DBDriver: os.Getenv("DB_DRIVER"),
	}

	workers := forge.ForgeWorkers{
		stub.ForgeAppGo,
		stub.ForgeBootstrap,
		stub.ForgeBootstrapSession,
		stub.ForgeMiddlewareAuth,
		stub.ForgeRouteHomeRegister,
	}

	if dbType.Name == "ent" {
		workers = append(workers,
			stub.CopyFile("ent.schema.user.stub", "ent/schema/user.go"),
			stub.CopyFile("ent.generate.stub", "ent/generate.go"),
			stub.CopyFile("ent.user.wrapper.stub", "internal/entauth/user.go"),
			stub.CopyFile("ent.userprovider.stub", "internal/entauth/entuserprovider.go"),
		)
	}

	driver := args.DBDriver
	if driver == "sqlite3" {
		driver = "sqlite"
	}

	switch driver {
	case "sqlite", "mysql", "postgres":
		for _, filename := range migrations {
			workers = append(
				workers,
				stub.CopyFile(
					fmt.Sprintf("migrations.%s/%s", driver, filename),
					fmt.Sprintf("database/migrations/%s", filename),
				),
			)
		}
	}

	files, forges, err := workers.Ready(ctx, args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	sort.Strings(files)

	return files, forges
}

func reloadDotEnv(dir string) {
	if err := cli.ReloadDotEnv(dir); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func askDBType() db.DBType {
	dbType, err := question.AskDBType()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dbType
}

func configureDotEnv(dir string) {
	err := run(
		dir,

		question.AskAppName,
		generateNewAppKey,
		question.AskDBDriver,
	)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createDotEnvFileAndLoad(dir string) {
	if err := cli.CreateDotEnvFileAndLoad(dir); errors.Is(err, cli.ErrOverwriteRejected) {
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func getProjectDir() string {
	dir := strings.TrimSpace(flag.Arg(0))
	if dir == "" {
		fmt.Println("Usage: golava-new <project>")
		os.Exit(1)
	}
	return dir
}

func run(dir string, processes ...func(dir string) error) error {
	for _, process := range processes {
		if err := process(dir); err != nil {
			return err
		}
	}
	return nil
}

func generateNewAppKey(dir string) error {
	fmt.Println("Generating new key...")
	key, err := encryption.GenerateKey()
	if err != nil {
		return err
	}

	err = cli.SetKeyInEnvironmentFile(dir, "APP_KEY", fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString(key)))
	if err != nil {
		return err
	}

	fmt.Println("Key generated successfully")
	return nil
}

func gitClone(remote string, project string) error {
	cmd := exec.Command("git", "clone", "--depth=1", "--branch", VERSION, remote, project)
	return cmd.Run()
}

func runGoModTidy(dir string) {
	fmt.Println("Running go mod tidy...")

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Printf("go mod tidy failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func runGoGenerateEnt(dir string) {
	fmt.Println("Running go generate ./ent...")

	cmd := exec.Command("go", "generate", "./ent")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Printf("go generate ./ent failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func createProject(dir string) {
	zipFile := fmt.Sprintf("golava-%s.zip", VERSION)

	downloadFile(
		fmt.Sprintf("https://github.com/wolftotem4/golava/archive/refs/tags/%s.zip", VERSION),
		zipFile,
	)

	unzipGithubProject(dir, zipFile)

	deleteFile(zipFile)
}

func unzipGithubProject(dir string, zipFile string) {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		fmt.Printf("Open zip file failed: %s\n", err.Error())
		os.Exit(1)
	}
	defer archive.Close()

	pathPrefix := archive.File[0].Name

	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		src, err := file.Open()
		if err != nil {
			fmt.Printf("Open file failed: %s\n", err.Error())
			os.Exit(1)
		}
		defer src.Close()

		path := strings.TrimPrefix(file.Name, pathPrefix)

		err = os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), os.ModePerm)
		if err != nil {
			fmt.Printf("Create folder failed: %s\n", err.Error())
			os.Exit(1)
		}

		dst, err := os.Create(filepath.Join(dir, path))
		if err != nil {
			fmt.Printf("Create file failed: %s\n", err.Error())
			os.Exit(1)
		}
		defer dst.Close()

		_, err = dst.ReadFrom(src)
		if err != nil {
			fmt.Printf("Read file failed: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func createProjectFolder(dir string) {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Create project folder failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func downloadFile(url string, file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err := cli.Download(url, file)
		if err != nil {
			fmt.Printf("Download failed: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func deleteFile(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Printf("Delete file failed: %s\n", err.Error())
		os.Exit(1)
	}
}
