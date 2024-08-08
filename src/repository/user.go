package repository

import (
	"database/sql"
	"gym-api/src/interfaces"
	"gym-api/src/models"
)

type Users struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) interfaces.UserRepository {
	return &Users{db}
}

func (repo Users) Save(user models.User) (uint64, error) {
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

func (repo Users) Delete(ID uint64) error {
	statement, err := repo.db.Prepare("delete from users where id = $1")
	if err != nil {
		return nil
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}

	return nil
}

func (repo Users) FindByEmail(email string) (models.User, error) {
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

func (repo Users) FindByID(ID uint64) (models.User, error) {
	row, err := repo.db.Query("select id, name, email from users where id = $1", ID)
	if err != nil {
		return models.User{}, err
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if err = row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repo Users) FindPassword(userID uint64) (string, error) {
	row, err := repo.db.Query("select password from users where id = $1", userID)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if err = row.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (repo Users) UpdatePassword(userID uint64, password string) error {
	statement, err := repo.db.Prepare("update users set password = $1 where id = $2")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userID); err != nil {
		return err
	}

	return nil
}
