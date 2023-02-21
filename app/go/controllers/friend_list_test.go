package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"problem1/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	models.InitDbForTest()
	defer models.CloseDb()

	users := []models.User{
		{UserID: 1, Name: "user1"},
		{UserID: 2, Name: "user2"},
		{UserID: 3, Name: "user3"},
		{UserID: 4, Name: "user4"},
		{UserID: 5, Name: "user5"},
	}
	if err := models.CreateUsers(users); err != nil {
		panic(err)
	}
	defer func() {
		if err := models.DeleteAllUsers(); err != nil {
			panic(err)
		}
	}()

	m.Run()
}

func TestGetFriendList(t *testing.T) {

	type param struct{ userId string }
	type want struct {
		code    int
		body    string
		message string
	}

	tests := []struct {
		name    string
		setup   func() error
		cleanup func() error
		param   param
		want    want
	}{
		{
			name: "OK",
			setup: func() error {
				fls := []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 1, User2ID: 3},
				}
				return models.CreateFriendLinks(fls)
			},
			cleanup: func() error {
				return models.DeleteAllFriendLinks()
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":2,"Name":"user2"},{"UserID":3,"Name":"user3"}]` + "\n",
			},
		},
		{
			name:  "Friend Not Found",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name:  "UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name:  "UserId Empty",
			param: param{userId: ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
	}

	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatal(err)
				}
			}
			if tt.cleanup != nil {
				defer func() {
					if err := tt.cleanup(); err != nil {
						t.Fatal(err)
					}
				}()
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get_friend_list/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.param.userId)

			err := getFriendList(c)
			if err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.code, rec.Code)
				var expectedJson, actualJson []models.User
				json.Unmarshal([]byte(rec.Body.Bytes()), &expectedJson)
				json.Unmarshal([]byte(tt.want.body), &actualJson)
				assert.ElementsMatch(t, expectedJson, actualJson)
			} else {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.want.code, he.Code)
					assert.Equal(t, tt.want.message, he.Message)
				}
			}
		})
	}
}

func TestGetFriendListOfFriendList(t *testing.T) {

	type param struct{ userId string }
	type want struct {
		code    int
		body    string
		message string
	}

	tests := []struct {
		name    string
		setup   func() error
		cleanup func() error
		param   param
		want    want
	}{
		{
			name: "OK No Include Friend",
			setup: func() error {
				fls := []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 1, User2ID: 3},
					{User1ID: 3, User2ID: 4},
					{User1ID: 4, User2ID: 5},
				}
				return models.CreateFriendLinks(fls)
			},
			cleanup: func() error {
				return models.DeleteAllFriendLinks()
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK Include Friend",
			setup: func() error {
				fls := []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 1, User2ID: 3},
					{User1ID: 2, User2ID: 3},
					{User1ID: 3, User2ID: 4},
					{User1ID: 4, User2ID: 5},
				}
				return models.CreateFriendLinks(fls)
			},
			cleanup: func() error {
				return models.DeleteAllFriendLinks()
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":2,"Name":"user2"},{"UserID":3,"Name":"user3"},{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name:  "Friend Not Found",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "Friend of Friend Not Found",
			setup: func() error {
				fls := []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 1, User2ID: 3},
				}
				return models.CreateFriendLinks(fls)
			},
			cleanup: func() error {
				return models.DeleteAllFriendLinks()
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name:  "UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name:  "UserId Empty",
			param: param{userId: ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
	}

	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatal(err)
				}
			}
			if tt.cleanup != nil {
				defer func() {
					if err := tt.cleanup(); err != nil {
						t.Fatal(err)
					}
				}()
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get_friend_list_of_friend_list/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.param.userId)

			err := getFriendListOfFriendList(c)
			if err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.code, rec.Code)
				var expectedJson, actualJson []models.User
				json.Unmarshal([]byte(rec.Body.Bytes()), &expectedJson)
				json.Unmarshal([]byte(tt.want.body), &actualJson)
				assert.ElementsMatch(t, expectedJson, actualJson)
			} else {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.want.code, he.Code)
					assert.Equal(t, tt.want.message, he.Message)
				}
			}
		})
	}
}
