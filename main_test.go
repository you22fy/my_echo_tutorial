package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

var initialUsers = map[string]User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

// setupEchoではEchoのインスタンスを初期化し、ルートエンドポイントの登録をする
func setUpEcho() *echo.Echo {
	users = make(map[string]User)
	for k, v := range initialUsers {
		users[k] = v
	}

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "password" {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/hello", helloHandler)
	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.POST("/users", createUserHandler)
	e.PUT("/users/:id", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	return e
}

// helloのテスト
func TestHelloHandler(t *testing.T) {
	e := setUpEcho()

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, helloHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{""message": "Hello, World!"}`, rec.Body.String())
	}
}
