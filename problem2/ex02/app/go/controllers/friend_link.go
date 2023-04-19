package controllers

import (
	"net/http"
	"problem1/models"

	"github.com/labstack/echo/v4"
)

func createFriendLink(c echo.Context) error {
	type Params struct {
		User1ID *int `json:"user1_id"`
		User2ID *int `json:"user2_id"`
	}
	var params Params
	err := c.Bind(&params)
	if err != nil || params.User1ID == nil || params.User2ID == nil ||
		*params.User1ID == *params.User2ID ||
		*params.User1ID <= 0 || *params.User2ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}

	users, err := models.GetUsers([]int{*params.User1ID, *params.User2ID})
	if err != nil || len(users) != 2 {
		if len(users) == 0 || users[0].UserID != *params.User1ID {
			return echo.NewHTTPError(http.StatusNotFound, "user1_id is not found")
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "user2_id is not found")
		}
	}

	fl := &models.FriendLink{User1ID: *params.User1ID, User2ID: *params.User2ID}
	err = models.CreateFriendLink(fl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create friend link")
	}

	return c.JSON(http.StatusOK, fl)
}
