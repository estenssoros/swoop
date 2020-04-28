package config

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/estenssoros/swoop/resource"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	defaultLimit    = 200
	defaultTruncate = false
)

func (c *Config) setDefaults() {
	if c.WriteLimit == nil {
		c.WriteLimit = &defaultLimit
	}
	if c.Truncate == nil {
		c.Truncate = &defaultTruncate
	}
}

// Process process a config
func (c *Config) Process() error {
	c.setDefaults()
	var g errgroup.Group
	for _, t := range c.Tables {
		t := t
		g.Go(func() error {
			return c.processTable(t)
		})
	}
	return g.Wait()
}

func (c *Config) processTable(t *Table) error {
	logrus.Infof("processing %s -> %s", t.Source, t.Destination)
	destinationSchema, err := c.Destination.GetSchema(t.Destination)
	if err != nil {
		return errors.Wrap(err, "get destination schema")
	}
	rows, err := c.Source.Query(t.Source)
	if err != nil {
		return errors.Wrap(err, "query")
	}
	defer rows.Close()
	sourceSchema, err := resource.SchemaFromRows(c.Source.Flavor, rows)
	if err != nil {
		return errors.Wrap(err, "source columns")
	}
	if err := sourceSchema.MustEqual(destinationSchema); err != nil {
		return errors.Wrap(err, "must equal")
	}

	if *c.Truncate {
		if err := c.Destination.Truncate(t.Destination); err != nil {
			return errors.Wrap(err, "truncate")
		}
	}

	input := &tupleizeInput{
		len:               len(sourceSchema.Columns),
		rows:              rows,
		destinationSchema: destinationSchema,
		sourceSchema:      sourceSchema,
	}

	writeInput := &writeTuplesInput{
		db:        c.Destination.DB,
		TableName: t.Destination,
		Schema:    destinationSchema,
		Tuples:    []string{},
	}
	var count int
	for tpl := range tupleize(input) {
		count++
		if err := tpl.err; err != nil {
			return err
		}
		writeInput.AddTuple(tpl.val)
		if writeInput.tupleLen() == *c.WriteLimit {
			if err := writeTuples(writeInput); err != nil {
				return errors.Wrap(err, "write tuples")
			}
			writeInput.ResetTuples()
		}
	}
	if writeInput.tupleLen() > 0 {
		if err := writeTuples(writeInput); err != nil {
			return errors.Wrap(err, "write tuples")
		}
	}
	logrus.Infof("wrote %d records to %s", count, t.Destination)
	return nil
}

func tupleizeRow(row []string) string {
	return fmt.Sprintf("(%s)", strings.Join(row, ","))
}

type tuple struct {
	val string
	err error
}

type tupleizeInput struct {
	len               int
	rows              *sql.Rows
	destinationSchema *resource.Schema
	sourceSchema      *resource.Schema
}

func (i *tupleizeInput) formatVal(idx int, col string) (string, error) {
	return i.destinationSchema.FormatSQLType(i.sourceSchema.Columns[idx].Name, col)
}

func tupleize(i *tupleizeInput) chan *tuple {
	ch := make(chan *tuple)
	values := make([]sql.RawBytes, i.len)
	scanArgs := make([]interface{}, i.len)
	row := make([]string, i.len)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	go func() {
		defer func() {
			close(ch)
		}()

		for i.rows.Next() {
			if err := i.rows.Scan(scanArgs...); err != nil {
				ch <- &tuple{err: err}
				return
			}

			for idx, col := range values {
				if col == nil {
					row[idx] = "NULL"
				} else {
					val, err := i.formatVal(idx, string(col))
					if err != nil {
						ch <- &tuple{err: errors.Wrap(err, "format sql type")}
						return
					}
					row[idx] = val
				}
			}
			ch <- &tuple{val: tupleizeRow(row)}
		}
	}()
	return ch
}
