package easypostgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	username string
	password string
	dbname   string
	DB       *sql.DB
}

func Open(username, password, dbname string) (*PostgreSQL, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, dbname))
	if err != nil {
		return nil, err
	}

	p := &PostgreSQL{
		username: username,
		password: password,
		dbname:   dbname,
		DB:       db,
	}

	return p, nil
}
