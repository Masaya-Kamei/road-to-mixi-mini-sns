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
		user_id, user_id, user_id,
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

func GetFriendListOfFriendListByUserId(user_id int) ([]User, error) {
	flFl := make([]User, 0)

	rows, err := db.Query(
		`select distinct u2.user_id, u2.name
		from users u1
		inner join friend_link fl1 on (u1.user_id = fl1.user1_id or u1.user_id = fl1.user2_id)
		inner join friend_link fl2 on (fl1.user1_id = fl2.user1_id or fl1.user1_id = fl2.user2_id or fl1.user2_id = fl2.user1_id or fl1.user2_id = fl2.user2_id)
		inner join users u2 on (u2.user_id = fl2.user1_id or u2.user_id = fl2.user2_id)
		where (fl1.user1_id = ? or fl1.user2_id = ?)
		and (fl2.user1_id != ? and fl2.user2_id != ?)
		and (u2.user_id != fl1.user1_id and u2.user_id != fl1.user2_id)
		`,
		user_id, user_id, user_id, user_id,
	);
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
			var fFl User
			err := rows.Scan(&fFl.UserID, &fFl.Name);
			if err != nil {
				return nil, err
			}
			flFl = append(flFl, fFl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flFl, nil
}
