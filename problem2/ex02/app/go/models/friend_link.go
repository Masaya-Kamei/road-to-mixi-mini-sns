package models

import "fmt"

type FriendLink struct {
	User1ID int
	User2ID int
}

func CreateFriendLink(fl *FriendLink) error {
	_, err := db.Exec(
		"insert into friend_link (user1_id, user2_id) values (?, ?)",
		fl.User1ID,
		fl.User2ID,
	)

	return err
}

func CreateFriendLinks(fls []FriendLink) error {
	query := "insert into friend_link (user1_id, user2_id) values "
	for _, fl := range fls {
		query += fmt.Sprintf("(%d, %d),", fl.User1ID, fl.User2ID)
	}
	_, err := db.Exec(query[:len(query)-1])

	return err
}

func GetFriendLinkByUserId(userID int) (FriendLink, error) {
	var fl FriendLink
	err := db.QueryRow(
		"select user1_id, user2_id from friend_link where user1_id = ? or user2_id = ?",
		userID,
		userID,
	).Scan(&fl.User1ID, &fl.User2ID)

	return fl, err
}

func GetAllFriendLinks() ([]FriendLink, error) {
	rows, err := db.Query("select user1_id, user2_id from friend_link")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fls := make([]FriendLink, 0)
	for rows.Next() {
		var fl FriendLink
		err := rows.Scan(&fl.User1ID, &fl.User2ID)
		if err != nil {
			return nil, err
		}
		fls = append(fls, fl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fls, nil
}

func DeleteAllFriendLinks() error {
	_, err := db.Exec("delete from friend_link")

	return err
}
