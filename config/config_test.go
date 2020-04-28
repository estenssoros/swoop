package config

import (
	"testing"
)

var testConfigString = `
source:
  name: theseus
  flavor: mysql
  secretPath: secret/data/theseus/database
destination:
  name: data-warehouse
  flavor: mssql
  secretPath: secret/data/data-warehouse/database
secretProvider:
  flavor: vault
  connectionURL: http://34.214.188.18:8200
`

func TestConfig(t *testing.T) {
	c, err := NewFromString(testConfigString)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}
}
