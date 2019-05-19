package main

import (
	"errors"
	"time"
)

type User struct {
	Id        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	State     int       `json:"state" db:"state"`
	Token     string    `json:"token" db:"token"`
}

func (u *User) CheckState() error {
	switch u.State {
	case 1:
		return errors.New("user not verified")
	case 2:
		return nil
	case 3:
		return errors.New("user changed password")
	case 4:
		return errors.New("user locked")
	default:
		return errors.New("user invalid state")
	}

}

type Login struct {
	UserId    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ExpireAt  time.Time `json:"expire_at" db:"expire_at"`
}
type NewUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type VerifyRegistration struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId string    `json:"user_id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

type LogoutRequest struct {
	UserId string `json:"user_id"`
}

type LogoutResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type PasswordResetRequest struct {
	Name string `json:"name"`
}
type PasswordResetResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

type PasswordChangeRequest struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type PasswordChangeResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
type LockRequest struct {
	UserID string `json:"user_id"`
}

type LockResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UnlockRequest struct {
	UserID string `json:"user_id"`
}
type UnlockResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
