package resource

import (
	"fmt"

	"github.com/pkg/errors"
)

type SQLType string

var (
	StringType   SQLType = "string"
	NumberType   SQLType = "number"
	DateType     SQLType = "date"
	DateTimeType SQLType = "datetime"
)

func (s SQLType) Format(val string) (string, error) {
	switch s {
	case StringType, DateType, DateTimeType:
		return fmt.Sprintf("'%s'", val), nil
	case NumberType:
		return val, nil
	default:
		return "", errors.Errorf("no type: %s", s)
	}
}

func ConvertTypeName(flavor Flavor, databaseTypeName string) (SQLType, error) {
	switch flavor {
	case MsSQLFlavor:
		return convertMsSQLDatabaseType(databaseTypeName)
	case MySQLFlavor:
		return convertMySQLDatabaseType(databaseTypeName)
	default:
		return "", errors.Wrap(ErrUnknownResourceFlavor, string(flavor))
	}
}

func convertMsSQLDatabaseType(databaseTypeName string) (SQLType, error) {
	switch databaseTypeName {
	case "INT", "BIT", "BIGINT", "FLOAT", "REAL", "DECIMAL", "DOUBLE":
		return NumberType, nil
	case "DATETIME", "DATETIME2":
		return DateTimeType, nil
	case "VARCHAR", "NVARCHAR", "NTEXT", "TEXT":
		return StringType, nil
	case "DATE":
		return DateType, nil
	default:
		return "", errors.Wrap(ErrUnknownDatabaseType, databaseTypeName)
	}
}
func convertMySQLDatabaseType(databaseTypeName string) (SQLType, error) {
	switch databaseTypeName {
	case "INT", "TINYINT", "BIGINT", "FLOAT", "DOUBLE", "DECIMAL":
		return NumberType, nil
	case "DATETIME":
		return DateTimeType, nil
	case "VARCHAR", "TEXT":
		return StringType, nil
	case "DATE":
		return DateType, nil
	default:
		return "", errors.Wrap(ErrUnknownDatabaseType, databaseTypeName)
	}
}
