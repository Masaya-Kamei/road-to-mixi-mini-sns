package models

import (
	"log"
	"problem1/database"
)

type User struct {
	UserID int
	Name   string
}

func (u *User) CreateUser() error {
	db := database.Get()

	_, err := db.Exec(
		"insert into users (user_id, name) values (?, ?)",
		u.UserID,
		u.Name,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetUserByUserId(user_id int) (User, error) {
	var user User
	db := database.Get()

	err := db.QueryRow(
		"select user_id, name from users where user_id = ?",
		user_id,
	).Scan(&user.UserID, &user.Name)
	if err != nil {
		log.Fatalln(err)
	}

	return user, err
}

func GetAllUsers() ([]User, error) {
	var users []User
	db := database.Get()

	rows, err := db.Query("select user_id, name from users")
	if err != nil {
			log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
			var user User
			err := rows.Scan(&user.UserID, &user.Name);
			if err != nil {
				log.Fatalln(err)
			}
			users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return users, nil
}
