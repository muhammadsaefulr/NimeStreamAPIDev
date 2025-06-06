package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
)

func CreateUserToUserModel(user *request.CreateUser) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func UpdateUserToUserModel(user *request.UpdateUser) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}
}

func UpdatePassOrVerifyToUserModel(user *request.UpdatePassOrVerify) *model.User {
	return &model.User{
		Password:      user.Password,
		VerifiedEmail: user.VerifiedEmail,
	}
}

func UserResponseToUserModel(user *model.User) *model.User {
	return &model.User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		VerifiedEmail: user.VerifiedEmail,
	}
}
