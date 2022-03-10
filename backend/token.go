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
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

func generateRefreshToken(user *User) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(jwtSecretKey))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &jwtCustomClaims{
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	//cookie.SameSite =http.SameSiteNoneMode
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(err error, c echo.Context) error {
	// Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
