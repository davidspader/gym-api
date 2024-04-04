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
		"INSERT INTO workouts (user_id, name) VALUES ($1, $2) RETURNING id",
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

func (repo Workouts) GetWorkoutsNamesByUserID(userID uint64) ([]models.Workout, error) {
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

func (repo Workouts) GetWorkoutNameByID(ID uint64) (models.Workout, error) {
	row, err := repo.db.Query(
		"SELECT * FROM workouts WHERE id = $1",
		ID,
	)
	if err != nil {
		return models.Workout{}, err
	}
	defer row.Close()

	var workout models.Workout

	if row.Next() {
		if err = row.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
		); err != nil {
			return models.Workout{}, err
		}
	}

	return workout, nil
}

func (repo Workouts) Update(workoutID uint64, userID uint64, workout models.Workout) error {
	statement, err := repo.db.Prepare(
		"UPDATE workouts SET name = $1 WHERE id = $2 AND user_id = $3",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(workout.Name, workoutID, userID); err != nil {
		return err
	}

	return nil
}

func (repo Workouts) Delete(workoutID uint64, userID uint64) error {
	statement, err := repo.db.Prepare(
		"DELETE FROM workouts WHERE id = $1 AND user_id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(workoutID, userID); err != nil {
		return err
	}

	return nil
}

func (repo Workouts) GetWorkoutByID(ID uint64, userID uint64) (models.Workout, error) {
	row, err := repo.db.Query(
		`
        SELECT w.id AS workout_id, w.user_id, w.name AS workout_name,
               e.id AS exercise_id, e.name AS exercise_name,
               e.weight, e.reps
        FROM workouts w
        LEFT JOIN exercises_workout ew ON w.id = ew.workout_id
        LEFT JOIN exercises e ON ew.exercise_id = e.id
        WHERE w.id = $1 AND w.user_id = $2
        `,
		ID,
		userID,
	)
	if err != nil {
		return models.Workout{}, err
	}
	defer row.Close()

	var workout models.Workout

	for row.Next() {
		var exerciseID sql.NullInt64
		var exerciseName sql.NullString
		var exerciseWeight sql.NullInt64
		var exerciseReps sql.NullInt64

		var exercise models.Exercise

		if err := row.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&exerciseID,
			&exerciseName,
			&exerciseWeight,
			&exerciseReps,
		); err != nil {
			return models.Workout{}, err
		}

		if exerciseID.Valid {
			exercise.ID = uint64(exerciseID.Int64)
		}

		if exerciseName.Valid {
			exercise.Name = exerciseName.String
		}

		if exerciseWeight.Valid {
			exercise.Weight = uint16(exerciseWeight.Int64)
		}

		if exerciseReps.Valid {
			exercise.Reps = uint16(exerciseReps.Int64)
		}

		workout.Exercises = append(workout.Exercises, exercise)
	}

	return workout, nil
}

func (repo Workouts) AddExerciseToWorkout(workoutID uint64, exerciseID uint64) error {
	statement, err := repo.db.Prepare(
		"INSERT INTO exercises_workout (workout_id, exercise_id) VALUES ($1, $2)",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(workoutID, exerciseID); err != nil {
		return err
	}

	return nil
}

func (repo Workouts) RemoveExerciseFromWorkout(workoutID uint64, exerciseID uint64) error {
	statement, err := repo.db.Prepare(
		"DELETE FROM exercises_workout WHERE workout_id = $1 AND exercise_id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(workoutID, exerciseID); err != nil {
		return err
	}

	return nil
}
