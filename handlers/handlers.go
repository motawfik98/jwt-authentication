package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username == "jon" && password == "shhh!" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{
		Name:  "Authorization",
		Value: t,
	})
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "Authorization",
		Expires: time.Now(),
	})
	return c.JSON(http.StatusTemporaryRedirect, map[string]string{
		"url": "/",
	})
}

func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func EnsureNotLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, _ := c.Cookie("Authorization")
		if cookie != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "You must be logged out",
			})
		}
		return next(c)
	}
}

//func EnsureLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		cookie, _ := c.Cookie("Authorization")
//		if cookie == nil {
//			return c.JSON(http.StatusBadRequest, map[string]string{})
//		}
//		return next(c)
//	}
//}
