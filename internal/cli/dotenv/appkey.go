package dotenv

import (
	"encoding/base64"
	"fmt"

	"github.com/wolftotem4/golava-core/encryption"
)

func GenerateNewAppKey(dir string) error {
	fmt.Println("Generating new key...")
	key, err := encryption.GenerateKey()
	if err != nil {
		return err
	}

	err = SetKeyInEnvironmentFile(dir, "APP_KEY", fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString(key)))
	if err != nil {
		return err
	}

	fmt.Println("Key generated successfully")
	return nil
}
