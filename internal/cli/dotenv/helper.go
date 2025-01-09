package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

const DotEnvFile = ".env"

func SetKeyInEnvironmentFile(dir, key, value string) error {
	file := filepath.Join(dir, DotEnvFile)

	content, err := os.ReadFile(file)
	if err != nil {
		return errors.WithStack(err)
	}

	content, appends := SetEnvVar(content, key, value, os.Getenv(key))
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

func SetEnvVar(content []byte, key, value, old string) (newcontent []byte, appends bool) {
	const space = `[\t\f\v ]*`
	const unit = space + "%s" + space
	const regEx = `(?m)^` + unit + `=` + unit + `$`

	env := string(content)

	r := regexp.MustCompile(fmt.Sprintf(regEx, regexp.QuoteMeta(key), regexp.QuoteMeta(old)))

	if r.MatchString(env) {
		newEnv := r.ReplaceAllString(env, fmt.Sprintf("%s=%s", key, value))
		return []byte(newEnv), false
	} else {
		return []byte(fmt.Sprintf("\n%s=%s\n", key, value)), true
	}
}
