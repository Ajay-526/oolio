package psql

import (
	"database/sql"
	"fmt"
	"oolio/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	env := config.Env()

	dataBaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", env.DBUSER, env.DBPASSWORD, env.DBHOST, env.DBPORT, env.DBNAME, env.DBSSLMODE)

	db, err := sql.Open("postgres", dataBaseUrl)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
