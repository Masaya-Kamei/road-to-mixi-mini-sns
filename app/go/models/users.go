package models

import (
	"log"
	"problem1/database"
)

type User struct {
	UserID int64
	Name   string
}

func (u *User) CreateUser() (err error) {
	cmd := `insert into users (user_id, name) values (?, ?)`

	_, err = database.Get().Exec(cmd, u.UserID, u.Name)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(user_id int) (user User, err error) {
	user = User{}
	cmd := `select user_id, name from users where user_id = ?`
	err = database.Get().QueryRow(cmd, user_id).Scan(
		&user.UserID,
		&user.Name,
	)
	if err != nil {
		log.Fatalln(err)
	}
	return user, err
}
