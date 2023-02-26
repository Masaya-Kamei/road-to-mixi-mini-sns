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

// bonus
func getFriendOfFriendListPagingWithCache(c echo.Context) error {
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

	flFl, foundRows, err := models.GetFriendListOfFriendListPagingWithCache(userID, limit, page)
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

// for debug
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
