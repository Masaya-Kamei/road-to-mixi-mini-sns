package controllers

import "github.com/labstack/echo/v4"

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", top)
	e.GET("/get_friend_list/:user_id", getFriendList)
	e.GET("/get_friend_list_of_friend_list", getFriendListOfFriendList)
	e.GET("/get_friend_of_friend_list_paging", getFriendOfFriendListPaging)

	return e
}
