package controller

import (
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService service.LoginService
	jwtService   service.JWTService
}

func NewLoginController(loginService service.LoginService, jwtService service.JWTService) *LoginController {
	return &LoginController{
		loginService: loginService,
		jwtService:   jwtService,
	}
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (controller *LoginController) Login(ctx *gin.Context) string {
	var credentials LoginCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		return ""
	}

	isAuthenticated := controller.loginService.Login(credentials.Username, credentials.Password)
	if isAuthenticated {
		return controller.jwtService.GenerateToken(credentials.Username, true)
	}
	return ""
}
