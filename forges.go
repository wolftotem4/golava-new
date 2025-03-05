package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/wolftotem4/golava-new/internal/cli/question"
	"github.com/wolftotem4/golava-new/internal/db"
	"github.com/wolftotem4/golava-new/internal/forge"
	"github.com/wolftotem4/golava-new/stub"
)

func prepareForges(ctx context.Context, dir string, dbType db.DBType) ([]string, []func(ctx context.Context) error) {
	args := forge.ForgeWorkArgs{
		Dir:      dir,
		DBType:   dbType,
		DBDriver: os.Getenv("DB_DRIVER"),
	}

	workers := forge.ForgeWorkers{
		stub.ForgeAppGo,
		stub.ForgeBootstrapApp,
		stub.ForgeBootstrapSession,
		stub.ForgeMiddlewareAuth,
		stub.ForgeRouteHomeRegister,
	}

	if dbType.Name == "ent" {
		workers = append(workers,
			stub.CopyFile("ent.schema.user.stub", "database/ent/schema/user.go"),
			stub.CopyFile("ent.generate.stub", "database/ent/generate.go"),
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

func doForges(ctx context.Context, jobs []func(ctx context.Context) error) {
	for _, job := range jobs {
		if err := job(ctx); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}
	}
}

func promptOverwriteConfirmation(files []string) {
	if !question.AskOverwrite(files) {
		fmt.Println("Operation aborted")
		os.Exit(1)
	}
}
