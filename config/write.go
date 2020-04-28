package config

import (
	"fmt"
	"strings"

	"github.com/estenssoros/swoop/resource"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type writeTuplesInput struct {
	db        *sqlx.DB
	TableName string
	Schema    *resource.Schema
	Tuples    []string
}

func (w writeTuplesInput) AddTuple(tpl string) {
	w.Tuples = append(w.Tuples, tpl)
}

func (w *writeTuplesInput) ResetTuples() {
	w.Tuples = nil
}

func (w *writeTuplesInput) tupleLen() int {
	return len(w.Tuples)
}

func (input *writeTuplesInput) craftInsertStmt() (s string) {
	defer func() {
		fmt.Println(s)
	}()
	columns := make([]string, len(input.Schema.Columns))
	for i := 0; i < len(input.Schema.Columns); i++ {
		columns[i] = input.Schema.Columns[i].Name
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES ", input.TableName, strings.Join(columns, ","))
}

func writeTuples(input *writeTuplesInput) error {
	if _, err := input.db.Exec(input.craftInsertStmt() + strings.Join(input.Tuples, ",")); err != nil {
		return errors.Wrap(err, "insert")
	}
	return nil
}
