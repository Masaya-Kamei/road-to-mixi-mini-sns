package controllers

import (
	"net/http"
	"problem1/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getFriendList(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUserByUserId(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	fl, err := models.GetFriendListByUserId(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get friend list")
	}

	return c.JSON(http.StatusOK, fl)
}

func getFriendListOfFriendList(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUserByUserId(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendListByUserId(userId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "failed to get friend list of friend list")
	}

	return c.JSON(http.StatusOK, flFl)
}

func getFriendOfFriendListPaging(c echo.Context) error {
	// FIXME
	return nil
}
