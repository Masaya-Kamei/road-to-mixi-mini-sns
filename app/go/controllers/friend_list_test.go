package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"problem1/models"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	models.InitDbForTest()
	defer models.CloseDb()

	fixtures, err := testfixtures.New(
		testfixtures.Database(models.GetDb()),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("../fixtures"),
	)
	if err != nil {
		panic(err)
	}

	if err := fixtures.Load(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestGetFriendList(t *testing.T) {

	type param struct { userId string }
	type want struct { code int; body string; message string}

	tests := []struct {
		name   string
		param  param
		want   want
	}{
		{
			name: "OK",
			param: param{userId: "1"},
			want: want{
				code: http.StatusOK,
				body: `[{"UserID":2,"Name":"user2"},{"UserID":3,"Name":"user3"},{"UserID":4,"Name":"user4"}]` + "\n",
			},
		},
		{
			name: "Length 0",
			param: param{userId: "5"},
			want: want{
				code: http.StatusOK,
				body: `[]`+ "\n",
			},
		},
		{
			name: "UserId Not Found",
			param: param{userId: "100"},
			want: want{
				code: http.StatusNotFound,
				message: `user_id is not found`,
			},
		},
		{
			name: "UserId Not Integer",
			param: param{userId: "a"},
			want: want{
				code: http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
		{
			name: "UserId Empty",
			param: param{userId: ""},
			want: want{
				code: http.StatusBadRequest,
				message: `user_id is not integer`,
			},
		},
	}

	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/get_friend_list/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.param.userId)

			err := getFriendList(c)
			if err == nil {
				assert.NoError(t, err);
				assert.Equal(t, tt.want.code, rec.Code)
				assert.Equal(t, tt.want.body, rec.Body.String())
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
