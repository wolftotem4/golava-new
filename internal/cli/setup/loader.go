package setup

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func LoadSetupConfig(file io.Reader) (SetupConfig, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return SetupConfig{}, errors.WithStack(err)
	}

	var config SetupConfig
	err = yaml.Unmarshal(content, &config)
	return config, errors.WithStack(err)
}
