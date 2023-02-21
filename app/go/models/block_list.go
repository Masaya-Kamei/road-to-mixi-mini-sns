package models

import "fmt"

type BlockList struct {
	User1ID int
	User2ID int
}

func CreateBlockList(bl *BlockList) error {
	_, err := db.Exec(
		"insert into block_list (user1_id, user2_id) values (?, ?)",
		bl.User1ID,
		bl.User2ID,
	)

	return err
}

func CreateBlockLists(bls []BlockList) error {
  query := "insert into friend_link (user1_id, user2_id) values "
	for _, bl := range bls {
		query += fmt.Sprintf("(%d, %d),", bl.User1ID, bl.User2ID)
	}
	_, err := db.Exec(query[:len(query)-1])

	return err
}

func GetBlockListByUserId(user_id int) (BlockList, error) {
	var bl BlockList

	err := db.QueryRow(
		"select user1_id, user2_id from block_list where user1_id = ?",
		user_id,
	).Scan(&bl.User1ID, &bl.User2ID)

	return bl, err
}

func GetAllBlockLists() ([]BlockList, error) {
	bls := make([]BlockList, 0)

	rows, err := db.Query("select user1_id, user2_id from block_list")
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	for rows.Next() {
			var bl BlockList
			err := rows.Scan(&bl.User1ID, &bl.User2ID);
			if err != nil {
				return nil, err
			}
			bls = append(bls, bl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bls, nil
}

func DeleteAllBlockLists() error {
	_, err := db.Exec("delete from block_list")
	return err
}
