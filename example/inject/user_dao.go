package inject

import (
	"context"
	inject "github.com/sipt/goinject"
)

func init() {
	inject.Register(func(injector inject.Injector) (err error) {
		var dao IUserDao = new(UserDao)
		err = injector.Register("service.IUserDao", dao)
		return
	})
}

type IUserDao interface{
	GetUser(ctx context.Context, id int64) (*User, error)
}
type UserDao struct{
	db *Engine `inject:"db.engine"`
}
func (u *UserDao) GetUser(ctx context.Context, id int64) (*User, error){
	return &User{}, nil
}
