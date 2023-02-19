package main

import (
	"problem1/configs"
	ctl "problem1/controllers"
	"problem1/database"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	database.Init()
	defer database.Close()

	e := echo.New()
	e.GET("/", ctl.Index)
	e.GET("/get_friend_list", ctl.GetFriendList)
	e.GET("/get_friend_list_of_friend_list", ctl.GetFriendListOfFriendList)
	e.GET("/get_friend_of_friend_list_paging", ctl.GetFriendOfFriendListPaging)

	conf := configs.Get()
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
