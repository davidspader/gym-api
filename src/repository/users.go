package repository

import (
	"database/sql"
	"gym-api/src/models"
)

type Users struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *Users {
	return &Users{db}
}

// func (repo Users) Create(user models.User) (uint64, error) {
// 	statement, err := repo.db.Prepare(
// 		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer statement.Close()

// 	result, err := statement.Exec(user.Name, user.Email, user.Password)
// 	if err != nil {
// 		return 0, err
// 	}

// 	lastInsertID, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return uint64(lastInsertID), nil
// }

func (repo Users) Create(user models.User) (uint64, error) {
	statement, err := repo.db.Prepare(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var userID uint64
	err = statement.QueryRow(user.Name, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
