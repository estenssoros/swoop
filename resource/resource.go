package resource

import (
	"database/sql"
	"fmt"

	"github.com/estenssoros/swoop/secret"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	_ "github.com/denisenkom/go-mssqldb" //mssql
	_ "github.com/go-sql-driver/mysql"   // mysql driver
)

type Resource struct {
	Name          string
	Flavor        Flavor
	SecretPath    string `yaml:"secretPath"`
	ConnectionURL string `yaml:"connectionURL"`

	DB *sqlx.DB
}

func (r *Resource) Connect(secretProvider *secret.Provider) error {
	if r.ConnectionURL == "" {
		secret, err := secretProvider.Get(r.SecretPath)
		if err != nil {
			return errors.Wrap(err, "get secret")
		}
		url, err := craftConnectionURL(r.Flavor, secret)
		if err != nil {
			return errors.Wrap(err, "craft connection url")
		}
		r.ConnectionURL = url
	}
	db, err := sqlx.Connect(string(r.Flavor), r.ConnectionURL)
	if err != nil {
		return errors.Wrap(err, "connect db")
	}
	db.SetMaxIdleConns(0)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "db ping")
	}
	r.DB = db
	logrus.Infof("connected to %s", r.Name)
	return nil
}

func (r *Resource) Query(tableName string) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM %s", tableName)
	color.Green(sql)
	return r.DB.Query(sql)
}

func craftConnectionURL(flavor Flavor, secret secret.Secret) (string, error) {
	switch flavor {
	case MySQLFlavor:
		return craftMySQLConnectionURL(secret)
	case MsSQLFlavor:
		return craftMsSQLConnectionURL(secret)
	default:
		return "", errors.Wrap(ErrUnknownResourceFlavor, "")
	}
}

func craftMySQLConnectionURL(s secret.Secret) (string, error) {
	var (
		user     string
		password string
		host     string
		database string
	)
	if err := s.SetString("user", &user); err != nil {
		return "", errors.Wrap(err, "user")
	}
	if err := s.SetString("password", &password); err != nil {
		return "", errors.Wrap(err, "password")
	}
	if err := s.SetString("host", &host); err != nil {
		return "", errors.Wrap(err, "host")
	}
	if err := s.SetString("database", &database); err != nil {
		return "", errors.Wrap(err, "database")
	}

	return fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", user, password, host, database), nil
}

func craftMsSQLConnectionURL(s secret.Secret) (string, error) {
	var (
		user     string
		password string
		host     string
		database string
	)
	if err := s.SetString("user", &user); err != nil {
		return "", errors.Wrap(err, "user")
	}
	if err := s.SetString("password", &password); err != nil {
		return "", errors.Wrap(err, "password")
	}
	if err := s.SetString("host", &host); err != nil {
		return "", errors.Wrap(err, "host")
	}
	if err := s.SetString("database", &database); err != nil {
		return "", errors.Wrap(err, "database")
	}
	return fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", user, password, host, database), nil
}

func (r *Resource) GetSchema(tableName string) (*Schema, error) {
	sql, err := r.CraftSelectOne(tableName)
	if err != nil {
		return nil, errors.Wrap(err, "craft select one")
	}
	color.Green(sql)
	rows, err := r.DB.Query(sql)
	if err != nil {
		return nil, errors.Wrap(err, "query rows")
	}
	defer rows.Close()
	return SchemaFromRows(r.Flavor, rows)
}

func (r *Resource) CraftSelectOne(tableName string) (string, error) {
	switch r.Flavor {
	case MsSQLFlavor:
		return craftMsSQLSelectOne(tableName), nil
	case MySQLFlavor:
		return craftMySQLSelectOne(tableName), nil
	default:
		return "", errors.Wrap(ErrUnknownResourceFlavor, string(r.Flavor))
	}
}

func craftMsSQLSelectOne(tableName string) string {
	return fmt.Sprintf("SELECT TOP 1 * FROM %s", tableName)
}

func craftMySQLSelectOne(tableName string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
}

func (r *Resource) Truncate(tableName string) error {
	sql := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	color.Red(sql)
	_, err := r.DB.Exec(sql)
	return err
}
