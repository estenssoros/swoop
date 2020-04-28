package secret

import "github.com/pkg/errors"

// ErrSecretDoesNotExists for when a secret doesnt exists
var ErrSecretDoesNotExists = errors.New("secret does not exists")

// ErrUnknownSecretFlavor for an unknown secret flavor
var ErrUnknownSecretFlavor = errors.New("unknown secret flavor")
