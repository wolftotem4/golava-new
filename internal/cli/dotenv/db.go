package dotenv

import (
	"github.com/wolftotem4/golava-new/internal/cli/setup"
)

func SetupDBDriver(dir string, config setup.SetupConfig) error {
	return SetKeyInEnvironmentFile(dir, "DB_DRIVER", config.DB.Driver)
}
