package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"problem1/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type NullString struct {
	String string
	Valid  bool
}

func TestMain(m *testing.M) {
	models.InitDbForTest()
	defer models.CloseDb()

	users := []models.User{
		{UserID: 1, Name: "user1"}, {UserID: 2,  Name: "user2"},
		{UserID: 3, Name: "user3"}, {UserID: 4,  Name: "user4"},
		{UserID: 5, Name: "user5"}, {UserID: 6,  Name: "user6"},
		{UserID: 7, Name: "user7"}, {UserID: 8,  Name: "user8"},
		{UserID: 9, Name: "user9"}, {UserID: 10, Name: "user10"},
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
			name: "OK: Basic",
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
			name:  "OK: Friend Not Found",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name:  "NotFound: UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "BadRequest: UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name:  "BadRequest: UserId Empty",
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
			name: "OK: No Include Friend",
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
			name: "OK: Include Friend",
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
				body: `[
					{"UserID":2,"Name":"user2"},{"UserID":3,"Name":"user3"},
					{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name:  "OK: Friend Not Found",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Friend of Friend Not Found",
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
			name:  "NotFound: UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "BadRequest: UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name:  "BadRequest: UserId Empty",
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
				json.Unmarshal([]byte(tt.want.body), &expectedJson)
				json.Unmarshal([]byte(rec.Body.Bytes()), &actualJson)
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
			name: "OK: No Include Friend",
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
			name: "OK: Include Friend",
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
			name: "OK: Include Blocked",
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
			name:  "OK: Friend Not Found",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Friend of Friend Not Found",
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
			name:  "NotFound: UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "BadRequest: UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name:  "BadRequest: UserId Empty",
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
				json.Unmarshal([]byte(tt.want.body), &expectedJson)
				json.Unmarshal([]byte(rec.Body.Bytes()), &actualJson)
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

func TestGetFriendOfFriendListPaging(t *testing.T) {

	type fixture struct{ fls []models.FriendLink; bls []models.BlockList }
	type param struct{
		userId string
		limit NullString
		page  NullString
	}
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
			name: "OK: Limit=undefined Page=undefined",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[
					{"UserID":3,"Name":"user3"},{"UserID":4, "Name":"user4"},
					{"UserID":5,"Name":"user5"},{"UserID":6, "Name":"user6"},
					{"UserID":7,"Name":"user7"},{"UserID":8, "Name":"user8"},
					{"UserID":9,"Name":"user9"},{"UserID":10,"Name":"user10"}]` + "\n",
			},
		},
		{
			name: "OK: Limit=2 Page=undefined",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "2", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":3,"Name":"user3"},{"UserID":4, "Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK: Limit=0 Page=undefined",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "0", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Limit=undefined Page=0",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				page: NullString{String: "0", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Limit=undefined Page=2",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				page: NullString{String: "2", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Limit=2 Page=0",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "2", Valid: true},
				page: NullString{String: "0", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name: "OK: Limit=2 Page=1",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "2", Valid: true},
				page: NullString{String: "1", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":3,"Name":"user3"},{"UserID":4, "Name":"user4"}]` + "\n",
			},
		},
		{
			name: "OK: Limit=2 Page=4",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "2", Valid: true},
				page: NullString{String: "4", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":9,"Name":"user9"},{"UserID":10, "Name":"user10"}]` + "\n",
			},
		},
		{
			name: "OK: Limit=2 Page=5",
			fixture: fixture{
				fls: []models.FriendLink{
					{User1ID: 1, User2ID: 2},
					{User1ID: 2, User2ID: 3}, {User1ID: 2, User2ID: 4},
					{User1ID: 2, User2ID: 5}, {User1ID: 2, User2ID: 6},
					{User1ID: 2, User2ID: 7}, {User1ID: 2, User2ID: 8},
					{User1ID: 2, User2ID: 9}, {User1ID: 2, User2ID: 10},
				},
			},
			param: param{
				userId: "1",
				limit: NullString{String: "2", Valid: true},
				page: NullString{String: "5", Valid: true},
			},
			want: want{
				code: http.StatusOK,
				body: `[]` + "\n",
			},
		},
		{
			name:  "NotFound: UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code:    http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name:  "BadRequest: UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code:    http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name:  "BadRequest: UserId Empty",
			param: param{userId: ""},
			want: want{
				code:    http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=-1 Page=undefined",
			param: param{
				userId: "1",
				limit: NullString{String: "-1", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=undefined Page=-1",
			param: param{
				userId: "1",
				page: NullString{String: "-1", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=ALL Page=undefined",
			param: param{
				userId: "1",
				limit: NullString{String: "ALL", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=undefined Page=a",
			param: param{
				userId: "1",
				page: NullString{String: "a", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=(empty) Page=undefined",
			param: param{
				userId: "1",
				limit: NullString{String: "", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
			},
		},
		{
			name: "BadRequest: Limit=undefined Page=(empty)",
			param: param{
				userId: "1",
				page: NullString{String: "", Valid: true},
			},
			want: want{
				code: http.StatusBadRequest,
				message: `invalid params`,
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

			q := make(url.Values)
			if tt.param.limit.Valid {
				q.Set("limit", tt.param.limit.String)
			}
			if tt.param.page.Valid {
				q.Set("page", tt.param.page.String)
			}
			var req *http.Request
			if len(q) == 0 {
				req = httptest.NewRequest(http.MethodGet, "/", nil)
			} else {
				req = httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get_friend_of_friend_list_of_paging/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.param.userId)

			err := getFriendOfFriendListPaging(c)
			if err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.code, rec.Code)
				var expectedJson, actualJson []models.User
				json.Unmarshal([]byte(tt.want.body), &expectedJson)
				json.Unmarshal([]byte(rec.Body.Bytes()), &actualJson)
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
