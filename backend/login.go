// Copyright (c) 2012-2022 GrabLink PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein is, and
// remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright
// law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code, which
// includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package main

import (
	"errors"
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
<<<<<<< HEAD
	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
=======
	//setTokenCookie(accessTokenCookieName, accessToken, exp, c)
>>>>>>> 8f1b52c (Initial commit)

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