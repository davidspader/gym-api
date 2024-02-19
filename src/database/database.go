package database

import (
	"database/sql"
	"gym-api/src/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	db, erro := sql.Open("postgres", config.DBString)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}
