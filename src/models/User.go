package models

import (
	"errors"
	"gym-api/src/security"
	"strings"

	"github.com/badoux/checkmail"
)

type User struct {
	ID       uint64 `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ChangePassword struct {
	New     string `json:"new"`
	Current string `json:"current"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	if err := user.format(step); err != nil {
		return err
	}

	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("name is required and cannot be blank")
	}

	if user.Email == "" {
		return errors.New("email is required and cannot be blank")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("the entered email is invalid")
	}

	if step == "register" && user.Password == "" {
		return errors.New("password is required and cannot be blank")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if step == "register" {
		hashPassword, err := security.GenerateHash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(hashPassword)
	}

	return nil
}
