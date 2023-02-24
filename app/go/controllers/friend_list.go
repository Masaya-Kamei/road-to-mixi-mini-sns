package controllers

import (
	"fmt"
	"net/http"
	"problem1/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func getFriendList(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	fl, err := models.GetFriendList(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get friend list")
	}

	return c.JSON(http.StatusOK, fl)
}

func getFriendListOfFriendList(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendList(userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "failed to get friend list of friend list")
	}

	return c.JSON(http.StatusOK, flFl)
}

func getFriendListOfFriendListExceptFriendAndBlocked(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, err := models.GetFriendListOfFriendListExceptFriendAndBlocked(userID)
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
		UserID *int `param:"user_id"`
		Limit *int `query:"limit"`
		Page *int `query:"page"`
	}

	var params Params = Params{UserID: nil, Limit: nil, Page: nil}
	err := c.Bind(&params)
	if ((err != nil) ||
		(params.UserID == nil || c.Param("user_id") == "") ||
		(params.Limit != nil && (c.QueryParam("limit") == "" || *params.Limit < 0)) ||
		(params.Page != nil && (c.QueryParam("page") == "" || *params.Page < 0))) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}
	userID, limit, page := *params.UserID, params.Limit, params.Page

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	flFl, foundRows, err := models.GetFriendListOfFriendListPaging(userID, limit, page)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend list of friend list paging",
		)
	}

	if (limit != nil && page != nil) {
		c.Response().Header().Set("Link", generateLinkHeader(c, *limit, *page, foundRows))
	}

	return c.JSON(http.StatusOK, flFl)
}

func generateLinkHeader(c echo.Context, limit, page, foundRows int) string {
	lastPage := (foundRows + 1) / limit
	baseUrl := "http://"+ c.Request().Host + c.Request().URL.Path
	linkHeaderTemplate := "<" + baseUrl + "?limit=%d&page=%d>; rel=\"%s\", "
	linkHeader := fmt.Sprintf(linkHeaderTemplate, limit, 1, "first")
	linkHeader += fmt.Sprintf(linkHeaderTemplate, limit, lastPage, "last")
	if (page > 1) {
		linkHeader += fmt.Sprintf(linkHeaderTemplate, limit, page - 1, "prev")
	}
	if (page < lastPage) {
		linkHeader += fmt.Sprintf(linkHeaderTemplate, limit, page + 1, "next")
	}
	return strings.TrimSuffix(linkHeader, ", ")
}
