package secret

import (
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// Provider provides the secret
type Provider struct {
	Flavor        Flavor
	ConnectionURL string `yaml:"connectionURL"`
}

// Get returns the secret interface
func (s *Provider) Get(secretPath string) (Secret, error) {
	switch s.Flavor {
	case VaultFlavor:
		return s.getVaultSecret(secretPath)
	case RawFlavor:
		return nil, errors.New("raw secrets must have connection url")
	default:
		return nil, errors.Wrap(ErrUnknownSecretFlavor, string(s.Flavor))
	}
}

func (s *Provider) getVaultSecret(secretPath string) (Secret, error) {
	token, err := getVaultToken()
	if err != nil {
		return nil, errors.Wrap(err, "get vault token")
	}
	client, err := api.NewClient(&api.Config{
		Address: s.ConnectionURL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "new client")
	}
	client.SetToken(token)
	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	mp, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.Wrap(err, "failed to parse vault response")
	}
	return &vaultSecret{mp}, nil
}
