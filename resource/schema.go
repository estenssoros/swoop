package resource

import (
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
)

type Schema struct {
	Columns []*Column
	Lookup  map[string]int `json:"-"`
}

func SchemaFromRows(flavor Flavor, rows *sql.Rows) (*Schema, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(err, "column types")
	}
	columns := make([]*Column, len(colTypes))
	lookup := map[string]int{}
	for i, c := range colTypes {
		sqlType, err := ConvertTypeName(flavor, c.DatabaseTypeName())
		if err != nil {
			return nil, errors.Wrapf(err, "convert database type name: %s", c.Name())
		}
		columns[i] = &Column{
			Name:    c.Name(),
			SQLType: sqlType,
		}
		lookup[c.Name()] = i
	}
	return &Schema{
		Columns: columns,
		Lookup:  lookup,
	}, nil
}

func (s *Schema) Equals(other *Schema) bool {
	if len(s.Columns) != len(other.Columns) {
		return false
	}
	for _, c := range s.Columns {
		idx, ok := other.Lookup[c.Name]
		if !ok {
			return false
		}
		if c.SQLType != other.Columns[idx].SQLType {
			return false
		}
	}
	return true
}

func (s *Schema) MustEqual(other *Schema) error {
	if len(s.Columns) != len(other.Columns) {
		return errors.Wrap(ErrSchemasNotEqual, "column len")
	}
	for _, c := range s.Columns {
		idx, ok := other.Lookup[c.Name]
		if !ok {
			return errors.Wrapf(ErrSchemasNotEqual, "other schema missing: %s", c.Name)
		}
		if c.SQLType != other.Columns[idx].SQLType {
			return errors.Wrapf(ErrSchemasNotEqual, "column type mismatch: %s source: %s, dest: %s", c.Name, c.SQLType, other.Columns[idx].SQLType)
		}
	}
	return nil
}

func (s Schema) String() string {
	ju, _ := json.MarshalIndent(s, "", " ")
	return string(ju)
}

func (s *Schema) FormatSQLType(colName, val string) (string, error) {
	idx, ok := s.Lookup[colName]
	if !ok {
		return "", errors.Errorf("could not find column: %s", colName)
	}
	return s.Columns[idx].SQLType.Format(val)
}
