package models

import (
	"errors"
	"strings"
)

type Workout struct {
	ID        uint64     `json:"id,omitempty"`
	UserID    uint64     `json:"userID,omitempty"`
	Name      string     `json:"name,omitempty"`
	Exercises []Exercise `json:"exercises,omitempty"`
}

func (workout *Workout) Prepare() error {
	if err := workout.validate(); err != nil {
		return err
	}

	workout.format()
	return nil
}

func (workout *Workout) validate() error {
	if workout.Name == "" {
		return errors.New("name is required and cannot be blank")
	}

	return nil
}

func (workout *Workout) format() {
	workout.Name = strings.TrimSpace(workout.Name)
}
