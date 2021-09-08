package goinject

import (
	"errors"
)

const InjectorTag = "inject"

var (
	ErrKeyDuplicate = errors.New("key duplicate")
	ErrNotFound = errors.New("key not found")
	shared = NewSingleInjector()
	funcs = make([]InjectFunc, 0, 16)
)
type InjectFunc func(Injector)error

func DefaultInjector() Injector {
	return shared
}

type Injector interface {
	Register(key interface{}, value interface{}) error
	Replace(key interface{}, value interface{})
	Get(key string) (interface{}, error)
	Remove(key string)
	InjectAll()
	Inject(v interface{})
	Clear()
}

func Register(f InjectFunc) {
	funcs  = append(funcs, f)
}

func InjectAll(injector Injector) (err error){
	for _, f := range funcs {
		err = f(injector)
		if err != nil {
			return err
		}
	}
	injector.InjectAll()
	return
}
