package users

// User models an app user
type User struct {
	ID        uint   `json:"id,omitempty" db:"id"`
	FirstName string `json:"firstName,omitempty" db:"first_name"`
	LastName  string `json:"lastName,omitempty" db:"last_name"`
	Username  string `json:"username,omitempty" db:"username"`
	Email     string `json:"email,omitempty" db:"email"`
	Password  string `json:"password,omitempty" db:"password"`
	Token     string `json:"token,omitempty" db:"-"`
}
