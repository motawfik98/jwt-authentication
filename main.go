package main

import (
	"./handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", handlers.Login, handlers.EnsureNotLoggedIn)
	e.GET("/logout", handlers.Logout)

	// Unauthenticated route
	e.GET("/", handlers.Accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:Authorization",
	}))

	middleware.ErrJWTMissing = echo.NewHTTPError(
		http.StatusBadRequest,
		map[string]string{
			"message": "You must be logged in",
		},
	)

	r.GET("", handlers.Restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
