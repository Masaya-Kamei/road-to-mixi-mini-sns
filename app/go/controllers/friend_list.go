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

func getFriendOfFriendList(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	ffl, err := models.GetFriendOfFriendList(userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "failed to get friend of friend list")
	}

	return c.JSON(http.StatusOK, ffl)
}

func getFriendOfFriendListExceptFriendAndBlocked(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is not integer")
	}

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	ffl, err := models.GetFriendOfFriendListExceptFriendAndBlocked(userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend of friend list except friend and blocked",
		)
	}

	return c.JSON(http.StatusOK, ffl)
}

func getFriendOfFriendListPaging(c echo.Context) error {
	type Params struct {
		UserID *int `param:"user_id"`
		Limit  *int `query:"limit"`
		Page   *int `query:"page"`
	}

	var params Params = Params{UserID: nil, Limit: nil, Page: nil}
	err := c.Bind(&params)
	if (err != nil) ||
		(params.UserID == nil || c.Param("user_id") == "") ||
		(params.Limit != nil && (c.QueryParam("limit") == "" || *params.Limit <= 0)) ||
		(params.Page != nil && (c.QueryParam("page") == "" || *params.Page <= 0)) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}
	userID, limit, page := *params.UserID, params.Limit, params.Page

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	ffl, foundRows, err := models.GetFriendOfFriendListPaging(userID, limit, page)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend of friend list paging",
		)
	}

	setLinkHeader(c, limit, page, foundRows)

	return c.JSON(http.StatusOK, ffl)
}

func setLinkHeader(c echo.Context, limit *int, page *int, foundRows int) {
	if limit == nil || page == nil || *limit <= 0 || *page <= 0 {
		return
	}
	limitNum := *limit
	pageNum := *page
	lastPageNum := (foundRows + 1) / limitNum
	baseUrl := "http://" + c.Request().Host + c.Request().URL.Path
	linkHeaderTemplate := "<" + baseUrl + "?limit=%d&page=%d>; rel=\"%s\", "
	linkHeader := fmt.Sprintf(linkHeaderTemplate, limitNum, 1, "first")
	linkHeader += fmt.Sprintf(linkHeaderTemplate, limitNum, lastPageNum, "last")
	if pageNum > 1 {
		linkHeader += fmt.Sprintf(linkHeaderTemplate, limitNum, pageNum-1, "prev")
	}
	if pageNum < lastPageNum {
		linkHeader += fmt.Sprintf(linkHeaderTemplate, limitNum, pageNum+1, "next")
	}
	c.Response().Header().Set("Link", strings.TrimSuffix(linkHeader, ", "))
}

// bonus
func getFriendOfFriendListPagingWithCache(c echo.Context) error {
	type Params struct {
		UserID *int `param:"user_id"`
		Limit  *int `query:"limit"`
		Page   *int `query:"page"`
	}

	var params Params = Params{UserID: nil, Limit: nil, Page: nil}
	err := c.Bind(&params)
	if (err != nil) ||
		(params.UserID == nil || c.Param("user_id") == "") ||
		(params.Limit != nil && (c.QueryParam("limit") == "" || *params.Limit <= 0)) ||
		(params.Page != nil && (c.QueryParam("page") == "" || *params.Page <= 0)) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}
	userID, limit, page := *params.UserID, params.Limit, params.Page

	_, err = models.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user_id is not found")
	}

	ffl, foundRows, err := models.GetFriendOfFriendListPagingWithCache(userID, limit, page)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to get friend of friend list paging with cache",
		)
	}

	setLinkHeader(c, limit, page, foundRows)

	return c.JSON(http.StatusOK, ffl)
}
