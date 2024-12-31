package dotenv

import (
	"github.com/wolftotem4/golava-new/internal/cli/setup"
)

func ConfigureDotEnv(dir string, config setup.SetupConfig) error {
	err := GenerateNewAppKey(dir)
	if err != nil {
		return err
	}

	err = SetupDBDriver(dir, config)
	return err
}
