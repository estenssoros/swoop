package secret

// Flavor different types of secrets
type Flavor string

var (
	// VaultFlavor for vault
	VaultFlavor Flavor = "vault"
	// RawFlavor for raw secret
	RawFlavor Flavor = "raw"
)
