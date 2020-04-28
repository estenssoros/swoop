package secret

import "github.com/pkg/errors"

type vaultSecret struct {
	secrets map[string]interface{}
}

func (s *vaultSecret) Get(key string) (interface{}, error) {
	val, ok := s.secrets[key]
	if !ok {
		return nil, errors.Wrap(ErrSecretDoesNotExists, key)
	}
	return val, nil
}

func (s *vaultSecret) GetString(key string) (string, error) {
	i, err := s.Get(key)
	if err != nil {
		return "", errors.Wrap(err, "get secret")
	}
	val, ok := i.(string)
	if !ok {
		return "", errors.New("converting interface to string")
	}
	return val, nil
}

func (s *vaultSecret) SetString(key string, valPtr *string) error {
	val, err := s.GetString(key)
	if err != nil {
		return errors.Wrap(err, "get string")
	}
	*valPtr = val
	return nil
}

func (s *vaultSecret) Has(key string) bool {
	_, ok := s.secrets[key]
	return ok
}

func (s *vaultSecret) MustGetString(key string) string {
	return s.secrets[key].(string)
}
