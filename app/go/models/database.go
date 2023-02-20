package models

import (
	"database/sql"
	"problem1/configs"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDb() {
	conf := configs.Get()

	var err error
	db, err = sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
}

func InitDbForTest() {
	conf := configs.Get()

	var err error
	db, err = sql.Open(conf.DB.Driver, conf.DB.DataSourceTest)
	if err != nil {
		panic(err)
	}
}

func GetDb() *sql.DB {
	return db
}

func CloseDb() {
	db.Close()
}
