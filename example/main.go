package main

import (
	"context"

	"github.com/sipt/goinject"
	"github.com/sipt/goinject/example/inject"
)

func main() {
	// init db
	engine := new(inject.Engine)
	goinject.Register(func(injector goinject.Injector) (err error) {
		err = injector.Register("db.engine", engine)
		return
	})
	/**
	{
		"db.Engine": *Engine
		"dao.IUserDao": *UserDao{ db: nil }
		"dao.IUserService": *UserService{ dao: nil }
	}
	*/
	injector := goinject.DefaultInjector()
	/**
	调用 InjectAll() 后，injector 内部相互引用注入
	{
		"db.Engine": *Engine
		"dao.IUserDao": *UserDao{ db: *Engine }
		"dao.IUserService": *UserService{ dao: *UserDao }
	}
	*/
	injector.InjectAll()

	handle := &ExampleHandle{}
	// 把值注入到 handle
	injector.Inject(handle)
	user, err := handle.GetUser(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	_ = user

	// 直接取值
	iface, err := injector.Get("dao.IUserDao")
	if err != nil {
		if err == goinject.ErrNotFound {
			panic("Not found")
		}
		panic(err)
	}
	dao := iface.(inject.IUserDao)
	user, err = dao.GetUser(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	_ = user

	// test 隔离
	testInjector := goinject.NewSingleInjector()
	testInjector.InjectAll()
	testHandle := &ExampleHandle{}
	testInjector.Inject(testHandle)
	if testHandle.service == nil || testHandle.service == handle.service {
		panic("inject failed")
	}
}

type ExampleHandle struct {
	service inject.IUserService `inject:"service.IUserService"`
}

func (e *ExampleHandle) GetUser(ctx context.Context, id int64) (*inject.User, error) {
	return e.service.GetUser(ctx, id)
}
