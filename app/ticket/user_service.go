package ticket

import "net/http"

type UserService struct {
}

func NewUserService(*http.Client) *UserService {
	return &UserService{}
}
