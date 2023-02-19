package database

import (
	"database/sql"
	"problem1/configs"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	conf := configs.Get()

	var err error
	db, err = sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
}

func Get() *sql.DB {
	return db
}

func Close() {
	db.Close()
}
