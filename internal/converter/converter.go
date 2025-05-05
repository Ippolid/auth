package converter

import (
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserInfoFromService преобразует UpdateRequest в UserInfo
func ToUserInfoFromService(req *auth_v1.UpdateRequest) *model.UserInfo {
	name := req.Info.Name
	email := req.Info.Email
	user := model.UserInfo{
		Name:  name,
		Email: email,
	}

	return &user
}

// ToDescFromAuthGet преобразует ListResponse в UserList
func ToDescFromAuthGet(req *model.User) *auth_v1.GetResponse {
	role := auth_v1.Role_USER // Значение по умолчанию
	if req.Role {
		role = auth_v1.Role_ADMIN // Если Role = true, устанавливаем ADMIN
	}
	return &auth_v1.GetResponse{
		User: &auth_v1.UserGet{
			Id: req.ID,
			Info: &auth_v1.UserInfo{
				Name:  req.User.Name,
				Email: req.User.Email,
			},
			Role:      role,
			CreatedAt: timestamppb.New(req.CreatedAt),
			UpdatedAt: timestamppb.New(req.CreatedAt),
		},
	}
}

// ToAuthCreateFromDesc преобразует CreateRequest в User
func ToAuthCreateFromDesc(req *auth_v1.CreateRequest) *model.User {
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
