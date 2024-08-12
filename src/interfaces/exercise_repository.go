package interfaces

import "gym-api/src/models"

type ExerciseRepository interface {
	Save(exercise models.Exercise) (uint64, error)
	FindByID(exerciseID uint64, userID uint64) (models.Exercise, error)
	FindByUserID(userID uint64) ([]models.Exercise, error)
	Update(exerciseID uint64, userID uint64, exercise models.Exercise) error
	Delete(exerciseID uint64, userID uint64) error
	VerifyOwnership(exerciseIDs []uint64, userID uint64) error
}
