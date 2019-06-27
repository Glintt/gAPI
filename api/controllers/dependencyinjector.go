package controllers

import (
	"github.com/Glintt/gAPI/api/users"
	userModels "github.com/Glintt/gAPI/api/users/models"
	routing "github.com/qiangxue/fasthttp-routing"
)

func getUserService(c *routing.Context) users.UserService {
	userService, _ := users.NewUserServiceWithUser(userModels.GetInternalAPIUser())
	return userService
}