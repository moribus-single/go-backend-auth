package models

type User struct {
	Username     string
	Password     string
	Refresh      string
	TokenCounter int
}
