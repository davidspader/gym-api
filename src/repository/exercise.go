package repository

import (
	"database/sql"
	"errors"
	"gym-api/src/interfaces"
	"gym-api/src/models"

	"github.com/lib/pq"
)

type Exercises struct {
	db *sql.DB
}

func NewExercisesRepository(db *sql.DB) interfaces.ExerciseRepository {
	return &Exercises{db}
}

func (repo Exercises) Save(exercise models.Exercise) (uint64, error) {
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

func (repo Exercises) FindByID(exerciseID uint64, userID uint64) (models.Exercise, error) {
	row, err := repo.db.Query(
		"select * from exercises where id = $1 and user_id = $2",
		exerciseID, userID,
	)
	if err != nil {
		return models.Exercise{}, err
	}
	defer row.Close()

	var exercise models.Exercise

	if row.Next() {
		if err = row.Scan(
			&exercise.ID,
			&exercise.UserID,
			&exercise.Name,
			&exercise.Weight,
			&exercise.Reps,
		); err != nil {
			return models.Exercise{}, err
		}
	}

	return exercise, nil
}

func (repo Exercises) FindByUserID(userID uint64) ([]models.Exercise, error) {
	rows, err := repo.db.Query(
		"SELECT * FROM exercises WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	for rows.Next() {
		var exercise models.Exercise

		if err = rows.Scan(
			&exercise.ID,
			&exercise.UserID,
			&exercise.Name,
			&exercise.Weight,
			&exercise.Reps,
		); err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (repo Exercises) Update(exerciseID uint64, userID uint64, exercise models.Exercise) error {
	statement, err := repo.db.Prepare(
		"UPDATE exercises SET name = $1, weight = $2, reps = $3 WHERE id = $4 AND user_id = $5",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(exercise.Name, exercise.Weight, exercise.Reps, exerciseID, userID); err != nil {
		return err
	}

	return nil
}

func (repo Exercises) Delete(exerciseID uint64, userID uint64) error {
	statement, err := repo.db.Prepare(
		"DELETE FROM exercises WHERE id = $1 AND user_id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(exerciseID, userID); err != nil {
		return err
	}

	return nil
}

func (repo Exercises) VerifyOwnership(exerciseIDs []uint64, userID uint64) error {
	query := `
		SELECT COUNT(*)
		FROM exercises
		WHERE id = ANY($1::bigint[]) AND user_id = $2
	`
	var count int
	err := repo.db.QueryRow(query, pq.Array(exerciseIDs), userID).Scan(&count)
	if err != nil {
		return err
	}

	if count != len(exerciseIDs) {
		return errors.New("one or more exercises do not belong to the user")
	}

	return nil
}
