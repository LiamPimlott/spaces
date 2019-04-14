package users

import "github.com/asaskevich/govalidator"

// User models an app user
type User struct {
	ID        uint   `json:"id,omitempty" db:"id"`
	FirstName string `json:"first_name,omitempty" db:"first_name" valid:"required~first_name is required"`
	LastName  string `json:"last_name,omitempty" db:"last_name" valid:"required~last_name is required"`
	Email     string `json:"email,omitempty" db:"email" valid:"required~email is required"`
	Password  string `json:"password,omitempty" db:"password" valid:"required~password is required"`
	Token     string `json:"token,omitempty" db:"-"`
}

// Valid validates a Space struct.
func (u User) Valid() (bool, error) { return govalidator.ValidateStruct(u) }
