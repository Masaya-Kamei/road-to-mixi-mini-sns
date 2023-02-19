package models

import (
	"log"
	"problem1/database"
)

type BlockList struct {
	User1ID int
	User2ID int
}

func (bl *BlockList) CreateBlockList() error {
	db := database.Get()

	_, err := db.Exec(
		"insert into block_list (user1_id, user2_id) values (?, ?)",
		bl.User1ID,
		bl.User2ID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetBlockListByUserId(user_id int) (BlockList, error) {
	var bl BlockList
	db := database.Get()

	err := db.QueryRow(
		"select user1_id, user2_id from block_list where user1_id = ?",
		user_id,
	).Scan(&bl.User1ID, &bl.User2ID)
	if err != nil {
		log.Fatalln(err)
	}

	return bl, err
}

func GetAllBlockLists() ([]BlockList, error) {
	var bls []BlockList
	db := database.Get()

	rows, err := db.Query("select user1_id, user2_id from block_list")
	if err != nil {
			log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
			var bl BlockList
			err := rows.Scan(&bl.User1ID, &bl.User2ID);
			if err != nil {
				log.Fatalln(err)
			}
			bls = append(bls, bl)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return bls, nil
}
