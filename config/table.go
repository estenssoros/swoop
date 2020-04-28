package config

import "encoding/json"

type Table struct {
	Source      string
	Destination string
	Query       string
}

func (t Table) String() string {
	ju, _ := json.MarshalIndent(t, "", " ")
	return string(ju)
}
