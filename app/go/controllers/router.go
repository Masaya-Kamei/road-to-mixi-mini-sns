package controllers

import "github.com/labstack/echo/v4"

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", top)
	e.GET("/get_friend_list/:user_id", getFriendList)
	// e.GET("/get_friend_list_of_friend_list/:user_id", getFriendListOfFriendList)
	e.GET("/get_friend_list_of_friend_list/:user_id", getFriendListOfFriendListExceptFriendAndBlocked)
	e.GET("/get_friend_of_friend_list_paging/:user_id", getFriendOfFriendListPaging)

	return e
}
