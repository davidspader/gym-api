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

func (repo Users) Update(ID uint64, user models.User) error {
	statement, err := repo.db.Prepare(
		"update users set name = $1, email = $2 where id = $3",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Email, ID); err != nil {
		return err
	}

	return nil
}

func (repo Users) GetByEmail(email string) (models.User, error) {
	row, err := repo.db.Query("select id, password from users where email = $1", email)
	if err != nil {
		return models.User{}, err
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if err = row.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}
