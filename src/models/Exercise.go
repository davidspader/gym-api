package models

import (
	"errors"
	"strings"
)

type Exercise struct {
	ID     uint64 `json:"id,omitempty"`
	UserID uint64 `json:"userID,omitempty"`
	Name   string `json:"name,omitempty"`
	Weight uint16 `json:"weight,omitempty"`
	Reps   uint16 `json:"reps,omitempty"`
}

type ExerciseIDs struct {
	IDs []uint64 `json:"ids"`
}

func (exercise *Exercise) Prepare() error {
	if err := exercise.validate(); err != nil {
		return err
	}

	exercise.format()
	return nil
}

func (exercise *Exercise) validate() error {
	if exercise.Name == "" {
		return errors.New("name is required and cannot be blank")
	}

	return nil
}

func (exercise *Exercise) format() {
	exercise.Name = strings.TrimSpace(exercise.Name)
}
