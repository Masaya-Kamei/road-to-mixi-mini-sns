package models

import (
	"fmt"
	"math"
	"strconv"

	"github.com/patrickmn/go-cache"
)

type User struct {
	UserID int
	Name   string
}

func CreateUser(u *User) error {
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

func GetUser(UserID int) (User, error) {
	var user User
	err := db.QueryRow(
		"select user_id, name from users where user_id = ?",
		UserID,
	).Scan(&user.UserID, &user.Name)

	return user, err
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query("select user_id, name from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Name)
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

func GetUsers(UserIDs []int) ([]User, error) {
	userIDsStr := ""
	for _, userID := range UserIDs {
		userIDsStr += strconv.Itoa(userID) + ","
	}
	userIDsStr = userIDsStr[:len(userIDsStr)-1]

	rows, err := db.Query(`
		select user_id, name from users
		where user_id in (` + userIDsStr + `)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Name)
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

func GetFriendList(userID int) ([]User, error) {
	rows, err := db.Query(`
		select u.user_id, u.name
		from users u
		inner join friend_link fl
		on (
			(u.user_id = fl.user1_id or u.user_id = fl.user2_id)
			and (fl.user1_id = ? or fl.user2_id = ?)
			and (u.user_id != ?)
		)
		`,
		userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fl := make([]User, 0)
	for rows.Next() {
		var f User
		err := rows.Scan(&f.UserID, &f.Name)
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

func GetFriendListOfFriendList(userID int) ([]User, error) {
	rows, err := db.Query(`
		select distinct u2.user_id, u2.name
		from users u1
		inner join friend_link fl1
		on (
			(u1.user_id = fl1.user1_id or u1.user_id = fl1.user2_id)
			and (fl1.user1_id = ? or fl1.user2_id = ?)
			and (u1.user_id != ?)
		)
		inner join friend_link fl2
		on (
			(u1.user_id = fl2.user1_id or u1.user_id = fl2.user2_id)
			and (fl2.user1_id != ? and fl2.user2_id != ?)
		)
		inner join users u2
		on (
			(fl2.user1_id = u2.user_id or fl2.user2_id = u2.user_id)
			and (u1.user_id != u2.user_id)
		)
		`,
		userID, userID, userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flFl := make([]User, 0)
	for rows.Next() {
		var fFl User
		err := rows.Scan(&fFl.UserID, &fFl.Name)
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

func GetFriendListOfFriendListExceptFriendAndBlocked(userID int) ([]User, error) {
	rows, err := db.Query(`
		select distinct u2.user_id, u2.name
		from users u1
		inner join friend_link fl1
		on (
			(u1.user_id = fl1.user1_id or u1.user_id = fl1.user2_id)
			and (fl1.user1_id = ? or fl1.user2_id = ?)
			and (u1.user_id != ?)
		)
		left join block_list bl
		on (
			(bl.user1_id = ? and bl.user2_id = u1.user_id)
		)
		inner join friend_link fl2
		on (
			(bl.id is null)
			and (u1.user_id = fl2.user1_id or u1.user_id = fl2.user2_id)
			and (fl2.user1_id != ? and fl2.user2_id != ?)
		)
		inner join users u2
		on (
			(fl2.user1_id = u2.user_id or fl2.user2_id = u2.user_id)
			and (u1.user_id != u2.user_id)
		)
		left join friend_link fl3
		on (
			(u2.user_id = fl3.user1_id or u2.user_id = fl3.user2_id)
			and (fl3.user1_id = ? or fl3.user2_id = ?)
		)
		where fl3.id is null
		`,
		userID, userID, userID, userID, userID, userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flFl := make([]User, 0)
	for rows.Next() {
		var fFl User
		err := rows.Scan(&fFl.UserID, &fFl.Name)
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

func GetFriendListOfFriendListPaging(userID int, limit, page *int) ([]User, int, error) {
	var limitNum, offset uint64 = math.MaxUint64, 0
	if limit != nil {
		limitNum = uint64(*limit)
	}
	if page != nil {
		if *page != 1 && (limitNum > math.MaxUint64 / uint64(*page - 1)) {
			offset = math.MaxUint64
		} else {
			offset = limitNum * uint64(*page - 1)
		}
	}

	rows, err := db.Query(`
		select sql_calc_found_rows distinct u2.user_id, u2.name
		from users u1
		inner join friend_link fl1
		on (
			(u1.user_id = fl1.user1_id or u1.user_id = fl1.user2_id)
			and (fl1.user1_id = ? or fl1.user2_id = ?)
			and (u1.user_id != ?)
		)
		left join block_list bl
		on (
			(bl.user1_id = ? and bl.user2_id = u1.user_id)
		)
		inner join friend_link fl2
		on (
			(bl.id is null)
			and (u1.user_id = fl2.user1_id or u1.user_id = fl2.user2_id)
			and (fl2.user1_id != ? and fl2.user2_id != ?)
		)
		inner join users u2
		on (
			(fl2.user1_id = u2.user_id or fl2.user2_id = u2.user_id)
			and (u1.user_id != u2.user_id)
		)
		left join friend_link fl3
		on (
			(u2.user_id = fl3.user1_id or u2.user_id = fl3.user2_id)
			and (fl3.user1_id = ? or fl3.user2_id = ?)
		)
		where fl3.id is null
		order by u2.user_id
		limit ? offset ?
		`,
		userID, userID, userID, userID, userID, userID, userID, userID,
		limitNum, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	flFl := make([]User, 0)
	for rows.Next() {
		var fFl User
		err := rows.Scan(&fFl.UserID, &fFl.Name)
		if err != nil {
			return nil, 0, err
		}
		flFl = append(flFl, fFl)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var foundRows int
	err = db.QueryRow("select found_rows()").Scan(&foundRows)
	if err != nil {
		return nil, 0, err
	}

	return flFl, foundRows, nil
}

// bonus
type FlFlPagingCache struct {
	flFl      []User
	foundRows int
}

func GetFriendListOfFriendListPagingWithCache(userID int, limit, page *int) ([]User, int, error) {
	cacheKey := "friend_of_friend_list_paging-" + strconv.Itoa(userID)
	if limit != nil {
		cacheKey += "-" + strconv.Itoa(*limit)
	}
	if page != nil {
		cacheKey += "-" + strconv.Itoa(*page)
	}

	flFlPagingCache, found := cacheInstance.Get(cacheKey)
	if found {
		flFlPagingCache := flFlPagingCache.(FlFlPagingCache)
		return flFlPagingCache.flFl, flFlPagingCache.foundRows, nil
	}

	flFl, foundRows, err := GetFriendListOfFriendListPaging(userID, limit, page)
	if err != nil {
		return nil, 0, err
	}

	newCache := FlFlPagingCache{flFl: flFl, foundRows: foundRows}
	cacheInstance.Set(cacheKey, newCache, cache.DefaultExpiration)

	return flFl, foundRows, nil
}
