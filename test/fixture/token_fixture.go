package fixture

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
	"github.com/muhammadsaefulr/NimeStreamAPI/test/helper"
)

var ExpiresAccessToken = time.Now().UTC().Add(time.Minute * time.Duration(config.JWTAccessExp))
var ExpiresRefreshToken = time.Now().UTC().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
var ExpiresResetPasswordToken = time.Now().UTC().Add(time.Minute * time.Duration(config.JWTResetPasswordExp))
var ExpiresVerifyEmailToken = time.Now().UTC().Add(time.Minute * time.Duration(config.JWTVerifyEmailExp))

func AccessToken(user *user_model.User) (string, error) {
	accessToken, err := helper.GenerateToken(user.ID.String(), ExpiresAccessToken, config.TokenTypeAccess)
	if err != nil {
		return accessToken, err
	}
	return accessToken, nil
}

func RefreshToken(user *user_model.User) (string, error) {
	refreshToken, err := helper.GenerateToken(user.ID.String(), ExpiresRefreshToken, config.TokenTypeRefresh)
	if err != nil {
		return refreshToken, err
	}
	return refreshToken, nil
}

func ResetPasswordToken(user *user_model.User) (string, error) {
	resetPasswordToken, err := helper.GenerateToken(
		user.ID.String(), ExpiresResetPasswordToken, config.TokenTypeResetPassword,
	)
	if err != nil {
		return resetPasswordToken, err
	}
	return resetPasswordToken, nil
}

func VerifyEmailToken(user *user_model.User) (string, error) {
	verifyEmailToken, err := helper.GenerateToken(user.ID.String(), ExpiresVerifyEmailToken, config.TokenTypeVerifyEmail)
	if err != nil {
		return verifyEmailToken, err
	}
	return verifyEmailToken, nil
}
