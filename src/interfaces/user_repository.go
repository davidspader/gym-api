package interfaces

import "gym-api/src/models"

type UserRepository interface {
	Save(user models.User) (uint64, error)
	Update(ID uint64, user models.User) error
	Delete(ID uint64) error
	FindByEmail(email string) (models.User, error)
	FindPassword(userID uint64) (string, error)
	UpdatePassword(userID uint64, password string) error
}
