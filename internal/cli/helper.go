package cli

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
)

const DotEnvFile = ".env"

func SetKeyInEnvironmentFile(dir, key, value string) error {
	file := path.Join(dir, DotEnvFile)

	content, err := os.ReadFile(file)
	if err != nil {
		return errors.WithStack(err)
	}

	content, appends := SetEnvVar(content, key, value)
	if appends {
		f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()

		_, err = f.Write(content)
		return errors.WithStack(err)
	} else {
		return errors.WithStack(os.WriteFile(file, content, 0644))
	}
}

func SetEnvVar(content []byte, key, value string) (newcontent []byte, appends bool) {
	env := string(content)
	old := os.Getenv(key)

	r := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\s*=\s*%s.*`, regexp.QuoteMeta(key), regexp.QuoteMeta(old)))

	if r.MatchString(env) {
		newEnv := r.ReplaceAllString(env, fmt.Sprintf("%s=%s", key, value))
		return []byte(newEnv), false
	} else {
		return []byte(fmt.Sprintf("\n%s=%s\n", key, value)), true
	}
}
