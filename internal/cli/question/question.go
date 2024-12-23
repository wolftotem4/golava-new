package question

import (
	"fmt"

	"github.com/wolftotem4/golava-new/internal/cli"
	"github.com/wolftotem4/golava-new/internal/db"
)

func AskAppName(dir string) error {
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

func AskDBType() (db.DBType, error) {
	fmt.Print("Enter the type of DB package you want to use (sqlx/gorm/ent): ")
	var dbTypeStr string
	fmt.Scanln(&dbTypeStr)

	dbType, ok := db.Data[dbTypeStr]
	if !ok {
		return db.DBType{}, fmt.Errorf("DB package %s not supported", dbType)
	}

	return dbType, nil
}

func AskDBDriver(dir string) error {
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

func AskOverwrite(files []string) bool {
	for _, file := range files {
		fmt.Println(file)
	}

	fmt.Print("Above files will be overwritten. Do you want to continue? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	return answer == "y"
}
