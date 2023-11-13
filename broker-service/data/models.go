package data

import "time"

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	Password          string    `json:"password"`
	Active            int       `json:"active"`
	RollId            int       `json:"roll_id"`
	Roll              string    `json:"roll"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}
