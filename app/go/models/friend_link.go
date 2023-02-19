package models

import (
	"log"
	"problem1/database"
)

type FriendLink struct {
	User1ID int
	User2ID int
}

func (fl *FriendLink) CreateFriendLink() error {
	db := database.Get()

	_, err := db.Exec(
		"insert into friend_link (user1_id, user2_id) values (?, ?)",
		fl.User1ID,
		fl.User2ID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetFriendLinkByUserId(user_id int) (FriendLink, error) {
	var fl FriendLink
	db := database.Get()

	err := db.QueryRow(
		"select user1_id, user2_id from friend_link where user1_id = ? or user2_id = ?",
		user_id,
		user_id,
	).Scan(&fl.User1ID, &fl.User2ID)
	if err != nil {
		log.Fatalln(err)
	}

	return fl, err
}

func GetAllFriendLinks() ([]FriendLink, error) {
	var fls []FriendLink
	db := database.Get()

	rows, err := db.Query("select user1_id, user2_id from friend_link")
	if err != nil {
			log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
			var fl FriendLink
			err := rows.Scan(&fl.User1ID, &fl.User2ID);
			if err != nil {
				log.Fatalln(err)
			}
			fls = append(fls, fl)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return fls, nil
}
