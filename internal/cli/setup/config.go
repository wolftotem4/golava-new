package setup

type SetupConfig struct {
	DB struct {
		Driver string `yaml:"driver"`
		Type   string `yaml:"type"`
	} `yaml:"db"`
}
