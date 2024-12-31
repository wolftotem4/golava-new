package dotenv

import "github.com/wolftotem4/golava-new/internal/cli/setup"

func SetupDotEnvFile(dir string, config setup.SetupConfig) error {
	err := CreateDotEnvFileAndLoad(dir)
	if err != nil {
		return err
	}

	err = ConfigureDotEnv(dir, config)
	if err != nil {
		return err
	}

	return ReloadDotEnv(dir)
}
