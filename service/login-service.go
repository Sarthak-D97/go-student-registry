package service

type LoginService struct {
	authorizedUser string
	password       string
}

func NewLoginService() *LoginService {
	return &LoginService{
		authorizedUser: "admin",
		password:       "password",
	}
}

func (ls *LoginService) Login(user, pass string) bool {
	return ls.authorizedUser == user && ls.password == pass
}
