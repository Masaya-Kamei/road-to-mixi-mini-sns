package controllers

import "github.com/labstack/echo/v4"

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", top)
	e.GET("/get_friend_list/:user_id", getFriendList)
	// e.GET("/get_friend_list_of_friend_list/:user_id", getFriendListOfFriendList)
	e.GET("/get_friend_list_of_friend_list/:user_id", getFriendListOfFriendListExceptFriendAndBlocked)
	e.GET("/get_friend_of_friend_list_paging/:user_id", getFriendOfFriendListPaging)

	// bonus
	e.GET("/get_friend_of_friend_list_paging_with_cache/:user_id", getFriendOfFriendListPagingWithCache)
	// for debug
	e.POST("/create_friend_link", createFriendLink)

	// client ip check
	e.GET("/get_client_ip", func (c echo.Context) error {
		return c.String(200, c.RealIP())
	})

	return e
}
