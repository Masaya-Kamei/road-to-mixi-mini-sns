package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", top)
	e.GET("/get_friend_list/:user_id", getFriendList)
	e.GET("/get_friend_of_friend_list/:user_id", getFriendOfFriendList)
	// e.GET("/get_friend_of_friend_list/:user_id", getFriendOfFriendListExceptFriendAndBlocked)
	e.GET("/get_friend_of_friend_list_v2/:user_id", getFriendOfFriendListExceptFriendAndBlocked)
	e.GET("/get_friend_of_friend_list_paging/:user_id", getFriendOfFriendListPaging)

	// bonus
	e.GET("/get_friend_of_friend_list_paging_with_cache/:user_id", getFriendOfFriendListPagingWithCache)
	// for debug
	e.POST("/create_friend_link", createFriendLink)

	// client ip test
	e.GET("/get_client_ip", func (c echo.Context) error {
		return c.String(200, c.RealIP())
	})

	return e
}
