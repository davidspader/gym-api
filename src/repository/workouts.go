package repository

import (
	"database/sql"
	"gym-api/src/models"
)

type Workouts struct {
	db *sql.DB
}

func NewWorkoutsRepository(db *sql.DB) *Workouts {
	return &Workouts{db}
}

func (repo Workouts) Create(workout models.Workout) (uint64, error) {
	statement, err := repo.db.Prepare(
		"INSERT INTO workouts (user_id, name) values ($1, $2) RETURNING id",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var workoutID uint64
	err = statement.QueryRow(workout.UserID, workout.Name).Scan(&workoutID)
	if err != nil {
		return 0, err
	}

	return workoutID, nil
}

func (repo Workouts) GetWorkoutsByUserID(userID uint64) ([]models.Workout, error) {
	rows, err := repo.db.Query(
		"SELECT * FROM workouts WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout

	for rows.Next() {
		var workout models.Workout

		if err = rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
		); err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	return workouts, nil
}
