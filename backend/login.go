package main

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"jwt-http-only/backend/vendor/github.com/labstack/echo/v4"
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
	//setTokenCookie(accessTokenCookieName, accessToken, exp, c)

	refreshToken, exp, err := generateRefreshToken(u)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
		"access_token": accessToken,
	})
}


func UserAPI(ctx echo.Context) error {
	accessToken := ctx.Request().Header.Get("Authorization")

	token, err := jwt.ParseWithClaims(accessToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		log.Errorf("err is : %s", err)
		return err
	}

	claims := token.Claims.(*jwtCustomClaims)
	log.Infof("claims.Name : %+v", claims.Name)

	if !token.Valid {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthorized",
		})
	}
	return ctx.JSON(http.StatusOK, echo.Map{
		"name": claims.Name,
	})
}


func RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return c.NoContent(http.StatusNoContent)
		}
		log.Errorf("err is : %s", err)
		return err
	}
	log.Infof("cookie value : %+v", cookie.Value)

	token, err := jwt.ParseWithClaims(cookie.Value, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetRefreshJWTSecret()), nil
	})
	if err != nil {
		log.Errorf("err is : %s", err)
		return err
	}

	if !token.Valid {
		return errors.New("error")
	}

	claims := token.Claims.(*jwtCustomClaims)
	log.Infof("claims.Name : %+v", claims.Name)

	u := &User{Name: claims.Name}

	accessToken, exp, err := generateAccessToken(u)
	if err != nil {
		return err
	}

	rToken, exp, err := generateRefreshToken(u)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, rToken, exp, c)

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
	})
}