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
		{UserID: 1, Name: "user1"}, {UserID: 2, Name: "user2"},
		{UserID: 3, Name: "user3"}, {UserID: 4, Name: "user4"},
		{UserID: 5, Name: "user5"},
	}
	err1 := models.DeleteAllUsers()
	err2 := models.DeleteAllFriendLinks()
	err3 := models.DeleteAllBlockLists()
	err4 := models.CreateUsers(users)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("setup failed")
	}
	defer func() {
		err := models.DeleteAllUsers()
		if err != nil {
			panic("cleanup failed")
		}
	}()

	m.Run()
}

func TestGetFriendList(t *testing.T) {

	type fixture struct{ fls []models.FriendLink }
	type param struct{ userId string }
	type want struct {
		code    int
		body    string
		message string
	}

	tests := []struct {
		name    string
		fixture fixture
		param   param
		want    want
	}{
		{
			name: "OK",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
				},
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
			if tt.fixture.fls != nil {
				err := models.CreateFriendLinks(tt.fixture.fls)
				if err != nil {
					t.Fatal("setup failed")
				}
			}
			defer func() {
				err := models.DeleteAllFriendLinks()
				if err != nil {
					t.Fatal("cleanup failed")
				}
			}()

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

	type fixture struct{ fls []models.FriendLink }
	type param struct{ userId string }
	type want struct {
		code    int
		body    string
		message string
	}

	tests := []struct {
		name    string
		fixture fixture
		param   param
		want    want
	}{
		{
			name: "OK No Include Friend",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
					{User1ID: 3, User2ID: 4}, {User1ID: 4, User2ID: 5},
				},
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK Include Friend",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
					{User1ID: 2, User2ID: 3}, {User1ID: 3, User2ID: 4},
					{User1ID: 4, User2ID: 5},
				},
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
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
				},
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
			if tt.fixture.fls != nil {
				err := models.CreateFriendLinks(tt.fixture.fls)
				if err != nil {
					t.Fatal("setup failed")
				}
			}
			defer func() {
				err := models.DeleteAllFriendLinks()
				if err != nil {
					t.Fatal("cleanup failed")
				}
			}()

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


func TestGetFriendListOfFriendListExceptFriendAndFriendBlocked(t *testing.T) {

	type fixture struct{ fls []models.FriendLink; bls []models.BlockList }
	type param struct{ userId string }
	type want struct {
		code    int
		body    string
		message string
	}

	tests := []struct {
		name    string
		fixture fixture
		param   param
		want    want
	}{
		{
			name: "OK No Include Friend",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
					{User1ID: 3, User2ID: 4}, {User1ID: 4, User2ID: 5},
				},
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK Include Friend",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
					{User1ID: 2, User2ID: 3}, {User1ID: 3, User2ID: 4},
					{User1ID: 4, User2ID: 5},
				},
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK Include Blocked",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
					{User1ID: 2, User2ID: 4}, {User1ID: 3, User2ID: 5},
					{User1ID: 4, User2ID: 5},
				},
				bls: []models.BlockList{
					{User1ID: 1, User2ID: 2},
				},
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":5,"Name":"user5"}]` + "\n",
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
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2}, {User1ID: 1, User2ID: 3},
				},
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
		if tt.fixture.fls != nil {
			err := models.CreateFriendLinks(tt.fixture.fls)
			if err != nil {
				t.Fatal("setup failed")
			}
		}
		if tt.fixture.bls != nil {
			err := models.CreateBlockLists(tt.fixture.bls)
			if err != nil {
				t.Fatal("setup failed")
			}
		}
		defer func() {
			err1 := models.DeleteAllFriendLinks()
			err2 := models.DeleteAllBlockLists()
			if err1 != nil || err2 != nil {
				t.Fatal("cleanup failed")
			}
		}()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get_friend_list_of_friend_list/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.param.userId)

			err := getFriendListOfFriendListExceptFriendAndBlocked(c)
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
