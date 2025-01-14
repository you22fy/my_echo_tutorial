package main

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// サンプルユーザーのデータを保持するマップ

var users = map[string]User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

func main() {
	// echo instance
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// basic 認証
		if username == "admin" && password == "password" {
			return true, nil
		}
		return false, nil
	}))

	// エンドポイントとそのハンドラ
	e.GET("/hello", helloHandler)
	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.POST("/users", createUserHandler)
	e.PUT("/users/:id", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)

	// 8080ポートでサーバーをスタート！
	e.Start(":8080")
}

// helloHandler
func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

// getUsersHandler
func getUsersHandler(c echo.Context) error {
	userList := []User{}

	for _, u := range users {
		userList = append(userList, u)
	}

	// idソート
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].ID < userList[j].ID
	})

	return c.JSON(http.StatusOK, userList)
}

// getUserHandler
func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	u, ok := users[id]
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not fount",
		})
	}
	return c.JSON(http.StatusOK, u)
}

// createUserHandler
func createUserHandler(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}

	if u.ID == "" || u.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing fields",
		})
	}

	if _, exists := users[u.ID]; exists {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User already exists",
		})
	}

	users[u.ID] = *u

	return c.JSON(http.StatusCreated, u)
}

// updateUserHandler
func updateUserHandler(c echo.Context) error {
	id := c.Param("id")
	u, exists := users[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	updatedUser := new(User)
	if err := c.Bind(updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}

	if updatedUser.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing fields",
		})
	}

	u.Name = updatedUser.Name
	users[id] = u

	return c.JSON(http.StatusOK, u)
}

// deleteUserHandler
func deleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	_, exists := users[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	delete(users, id)

	return c.NoContent(http.StatusNoContent)
}
