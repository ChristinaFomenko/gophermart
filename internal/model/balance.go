package model

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
	UserID    int     `json:"-"`
}
