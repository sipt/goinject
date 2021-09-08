package inject

import (
	"context"
	inject "github.com/sipt/goinject"
)

func init() {
	inject.Register(func(injector inject.Injector) (err error) {
		var service IUserService = new(UserService)
		err = injector.Register("service.IUserService", service)
		return
	})
}

type IUserService interface{
	GetUser(ctx context.Context, id int64) (*User, error)
}
type UserService struct{
	dao IUserDao `inject:"dao.IUserDao"`
}

func (* UserService) GetUser(ctx context.Context, id int64) (*User, error) {
	return &User{}, nil
}