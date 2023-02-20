package models

type User struct {
	UserID int
	Name   string
}

func (u *User) CreateUser() error {

	_, err := db.Exec(
		"insert into users (user_id, name) values (?, ?)",
		u.UserID,
		u.Name,
	)

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
