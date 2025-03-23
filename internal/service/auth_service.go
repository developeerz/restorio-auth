package service

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (authService *AuthService) Refresh() {

}

func (authService *AuthService) CheckAccess() {

}
