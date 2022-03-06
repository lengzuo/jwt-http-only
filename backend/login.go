package main

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	jwt.StandardClaims
}


type User struct {
	Name string `json:"name"`
}

func Login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err !=nil {
		log.Errorf("error in bind: %s", err)
		return err
	}
	accessToken, exp, err := generateAccessToken(u)
	if err != nil {
		return err
	}
	setTokenCookie(accessTokenCookieName, accessToken, exp, c)

	refreshToken, exp, err := generateRefreshToken(u)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
	})
}


func UserAPI(ctx echo.Context) error {
	cookie, err := ctx.Cookie(accessTokenCookieName)
	if err != nil {
		log.Errorf("err is : %s", err)
		return err
	}
	log.Infof("cookie value : %+v", cookie.Value)

	token, err := jwt.ParseWithClaims(cookie.Value, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	claims := token.Claims.(*jwtCustomClaims)
	log.Infof("claims.Name : %+v", claims.Name)

	return ctx.JSON(http.StatusOK, echo.Map{
		"name": claims.Name,
	})
}