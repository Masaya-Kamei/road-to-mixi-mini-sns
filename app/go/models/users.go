package models

import (
	"fmt"
)

type User struct {
	UserID int
	Name   string
}

func CreateUser(u *User)  error {

	_, err := db.Exec(
		"insert into users (user_id, name) values (?, ?)",
		u.UserID,
		u.Name,
	)

	return err
}

func CreateUsers(us []User) error {
  query := "insert into users (user_id, name) values "
	for _, u := range us {
		query += fmt.Sprintf("(%d, '%s'),", u.UserID, u.Name)
	}
	_, err := db.Exec(query[:len(query)-1])

	return err
}

func GetUserByUserId(user_id int) (User, error) {
	var user User

	err := db.QueryRow(
		"select user_id, name from users where user_id = ?",
		user_id,
	).Scan(&user.UserID, &user.Name)

	return user, err
}

func GetAllUsers() ([]User, error) {
	users := make([]User, 0)

	rows, err := db.Query("select user_id, name from users")
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	for rows.Next() {
			var user User
			err := rows.Scan(&user.UserID, &user.Name);
			if err != nil {
				return nil, err
			}
			users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteAllUsers() error {
	_, err := db.Exec("delete from users")
	return err
}

func GetFriendListByUserId(user_id int) ([]User, error) {
	fl := make([]User, 0)

	rows, err := db.Query(
		`select u.user_id, u.name
		from users u
		inner join friend_link fl on (u.user_id = fl.user1_id or u.user_id = fl.user2_id)
		where (fl.user1_id = ? or fl.user2_id = ?)
		and u.user_id != ?`,
		user_id,
		user_id,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
			var f User
			err := rows.Scan(&f.UserID, &f.Name);
			if err != nil {
				return nil, err
			}
			fl = append(fl, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fl, nil
}
