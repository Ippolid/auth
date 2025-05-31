package converter

import (
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserInfoFromService преобразует UpdateRequest в UserInfo
func ToUserInfoFromService(req *user_v1.UpdateRequest) *model.UserInfo {
	if req == nil || req.Info == nil {
		return nil
	}
	name := req.Info.Name
	email := req.Info.Email
	user := model.UserInfo{
		Name:  name,
		Email: email,
	}

	return &user
}

// ToDescFromUserGet преобразует ListResponse в UserList
func ToDescFromUserGet(req *model.User) *user_v1.GetResponse {
	if req == nil {
		return nil
	}
	role := user_v1.Role_USER // Значение по умолчанию
	if req.Role {
		role = user_v1.Role_ADMIN // Если Role = true, устанавливаем ADMIN
	}
	return &user_v1.GetResponse{
		User: &user_v1.UserGet{
			Id: req.ID,
			Info: &user_v1.UserInfo{
				Name:  req.User.Name,
				Email: req.User.Email,
			},
			Role:      role,
			CreatedAt: timestamppb.New(req.CreatedAt),
			UpdatedAt: timestamppb.New(req.CreatedAt),
		},
	}
}

// ToUserCreateFromDesc преобразует CreateRequest в User
func ToUserCreateFromDesc(req *user_v1.CreateRequest) *model.User {
	if req == nil || req.GetInfo() == nil || req.GetInfo().GetUser() == nil {
		return nil
	}

	name := req.GetInfo().GetUser().Name
	email := req.GetInfo().GetUser().Email
	password := req.GetInfo().GetPassword()
	// Преобразуем Role в bool (предполагая, что Role - это enum или int32)
	role := req.GetInfo().GetRole() > 0

	userInfo := model.UserInfo{
		Name:  name,
		Email: email,
	}

	user := model.User{
		User:     userInfo,
		Password: password,
		Role:     role,
	}
	return &user
}
