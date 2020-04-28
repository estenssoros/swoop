package secret

// Secret for all your secret needs
type Secret interface {
	Get(string) (interface{}, error)
	GetString(string) (string, error)
	SetString(string, *string) error
	Has(string) bool
	MustGetString(string) string
}
