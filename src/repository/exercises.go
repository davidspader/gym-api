package repository

import (
	"database/sql"
	"gym-api/src/models"
)

type Exercises struct {
	db *sql.DB
}

func NewExercisesRepository(db *sql.DB) *Exercises {
	return &Exercises{db}
}

func (repo Exercises) Create(exercise models.Exercise) (uint64, error) {
	statement, err := repo.db.Prepare(
		"INSERT INTO exercises (user_id, name, weight, reps) values ($1, $2, $3, $4) RETURNING id",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var exerciseID uint64
	err = statement.QueryRow(exercise.UserID, exercise.Name, exercise.Weight, exercise.Reps).Scan(&exerciseID)
	if err != nil {
		return 0, err
	}

	return exerciseID, nil
}