package interfaces

import "gym-api/src/models"

type WorkoutRepository interface {
	Save(workout models.Workout) (uint64, error)
	FindNamesByUserID(userID uint64) ([]models.Workout, error)
	FindNameByID(ID uint64) (models.Workout, error)
	FindByID(ID uint64, userID uint64) (models.Workout, error)
	Update(workoutID uint64, userID uint64, workout models.Workout) error
	Delete(workoutID uint64, userID uint64) error
	AddExercise(workoutID uint64, exerciseID uint64) error
	RemoveExercise(workoutID uint64, exerciseID uint64) error
}
