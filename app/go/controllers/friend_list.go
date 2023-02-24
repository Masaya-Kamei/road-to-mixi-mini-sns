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

	_, err = models.GetUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	fl, err := models.GetFriendList(userId)
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

	_, err = models.GetUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendList(userId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "failed to get friend list of friend list")
	}

	return c.JSON(http.StatusOK, flFl)
}

func getFriendListOfFriendListExceptFriendAndBlocked(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendListExceptFriendAndBlocked(userId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend list of friend list except friend and blocked",
		)
	}

	return c.JSON(http.StatusOK, flFl)
}

func getFriendOfFriendListPaging(c echo.Context) error {
	type Params struct {
		UserId *int `param:"user_id"`
		Limit *int `query:"limit"`
		Page *int `query:"page"`
	}

	var params Params = Params{UserId: nil, Limit: nil, Page: nil}
	err := c.Bind(&params)
	if ((err != nil) ||
		(params.UserId == nil || c.Param("user_id") == "") ||
		(params.Limit != nil && (c.QueryParam("limit") == "" || *params.Limit < 0)) ||
		(params.Page != nil && (c.QueryParam("page") == "" || *params.Page < 0))) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}
	userId, limit, page := *params.UserId, params.Limit, params.Page

	_, err = models.GetUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendListPaging(userId, limit, page)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend list of friend list paging",
		)
	}

	return c.JSON(http.StatusOK, flFl)
}
