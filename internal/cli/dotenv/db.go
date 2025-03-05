package dotenv

import (
	"github.com/wolftotem4/golava-new/internal/cli/setup"
)

func SetupDBSettings(dir string, config setup.SetupConfig) error {
	err := SetKeyInEnvironmentFile(dir, "DB_DRIVER", config.DB.Driver)
	if err != nil {
		return err
	}

	conn, ok := config.DB.Connections[config.DB.Driver]
	if ok {
		err = SetKeyInEnvironmentFile(dir, "DB_DSN", conn.DSN)
		if err != nil {
			return err
		}
	}

	return nil
}
