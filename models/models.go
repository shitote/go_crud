package models

type User struct {
	UserId int64 `json:"userid"`
	Name 	string 	`json:"name"`
	Age	int64	`json:"age"`
	Email	string	`json:"email"`
}