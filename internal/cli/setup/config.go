package setup

type SetupConfig struct {
	DB struct {
		Driver      string                      `yaml:"driver"`
		Type        string                      `yaml:"type"`
		Connections map[string]ConnectionConfig `yaml:"connections"`
	} `yaml:"db"`
}

type ConnectionConfig struct {
	DSN string `yaml:"dsn"`
}
