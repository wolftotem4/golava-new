package main

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava-new/internal/cli"
	"github.com/wolftotem4/golava-new/internal/db"
	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/stub"
)

const GIT_REMOTE = "https://github.com/wolftotem4/golava.git"
const VERSION = "v0.1.1"

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

	dir := strings.TrimSpace(flag.Arg(0))
	if dir == "" {
		fmt.Println("Usage: golava-new <project>")
		return
	}

	if err := gitClone(GIT_REMOTE, dir); err != nil {
		fmt.Printf("git clone failed: %s\n", err.Error())
		return
	}

	deleteGitRepo(dir)

	if err := cli.SetupDotEnvFile(dir); errors.Is(err, cli.ErrOverwriteRejected) {
		return
	} else if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	err := run(
		dir,

		askAppName,
		generateNewAppKey,
		askDBDriver,
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	dbType, err := askDBType()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// reload .env file
	if err := cli.ReloadDotEnv(dir); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	args := forge.ForgeWorkArgs{
		Dir:      dir,
		DBType:   dbType,
		DBDriver: os.Getenv("DB_DRIVER"),
	}

	files, forges, err := forge.ForgeWorkers{
		stub.ForgeAppGo,
		stub.ForgeBootstrapApp,
		stub.ForgeMiddlewareAuth,
		stub.ForgeRouteHomeRegister,
	}.Ready(ctx, args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	if !askForOverwrite(files) {
		fmt.Println("Operation aborted")
		return
	}

	for _, forge := range forges {
		if err := forge(ctx); err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
	}
}

func run(dir string, processes ...func(dir string) error) error {
	for _, process := range processes {
		if err := process(dir); err != nil {
			return err
		}
	}
	return nil
}

func askAppName(dir string) error {
	fmt.Print("Enter the name of your application: ")
	var appName string
	fmt.Scanln(&appName)

	err := cli.SetKeyInEnvironmentFile(dir, "APP_NAME", appName)
	if err != nil {
		return err
	}

	fmt.Println("Application name set successfully")
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

func askDBDriver(dir string) error {
	fmt.Print("Enter the driver of your DB (sqlite/mysql/postgres): ")
	var dbDriver string
	fmt.Scanln(&dbDriver)

	err := cli.SetKeyInEnvironmentFile(dir, "DB_DRIVER", dbDriver)
	if err != nil {
		return err
	}

	fmt.Println("DB driver set successfully")
	return nil
}

func askDBType() (db.DBType, error) {
	fmt.Print("Enter the type of DB package you want to use (sqlx/gorm): ")
	var dbTypeStr string
	fmt.Scanln(&dbTypeStr)

	dbType, ok := db.Data[dbTypeStr]
	if !ok {
		return db.DBType{}, fmt.Errorf("DB package %s not supported", dbType)
	}

	return dbType, nil
}

func askForOverwrite(files []string) bool {
	for _, file := range files {
		fmt.Println(file)
	}

	fmt.Print("Above files will be overwritten. Do you want to continue? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	return answer == "y"
}

func gitClone(remote string, project string) error {
	cmd := exec.Command("git", "clone", "--depth=1", "--branch", VERSION, remote, project)
	return cmd.Run()
}

func deleteGitRepo(dir string) error {
	return os.RemoveAll(path.Join(dir, "/.git"))
}
