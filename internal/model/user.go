package model

type User struct {
	ID       int    `json:"id,omitempty" db:"user_id"`
	Login    string `json:"login,omitempty" db:"login"`
	Password string `json:"password,omitempty" db:"password"`
}
