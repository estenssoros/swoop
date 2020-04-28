package resource

import "errors"

var ErrUnknownResourceFlavor = errors.New("unknown flavor flavor")
var ErrUnknownDatabaseType = errors.New("unknown database type")
var ErrSchemasNotEqual = errors.New("schemas not equal")
