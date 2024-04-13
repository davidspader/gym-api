package usecases

import (
	"gym-api/src/database"
	"gym-api/src/interfaces"
	"gym-api/src/models"
	"gym-api/src/repository"
)

type UserUseCase struct {
	repo interfaces.UserRepository
}

func NewUsersUseCase(repo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (c *UserUseCase) Save(user models.User) (uint64, error) {
	db, err := database.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	user.ID, err = repo.Save(user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
