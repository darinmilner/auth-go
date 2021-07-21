package models

import "time"

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	UserName          string    `json:"username"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	DeletedAt         time.Time `json:"-"`
	PasswordResetCode string    `json:"passwordResetCode,omitempty"`
}

type ForgotPasswordEmailPayload struct {
	Source            string `json:"source"`
	Destination       string `json:"destination"`
	PasswordResetCode string `json:"passwordResetCode,omitempty"`
}
