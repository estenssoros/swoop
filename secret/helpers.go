package secret

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// check environment for a variable
func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return val, fmt.Errorf("missing environment variable: %s", key)
	}
	return val, nil
}

func getVaultToken() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	file, err := ioutil.ReadFile(path.Join(homeDir, ".vault-token"))
	if err != nil {
		token, err := getEnv("VAULT_TOKEN")
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return string(file), nil
}
