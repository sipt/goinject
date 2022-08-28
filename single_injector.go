package goinject

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var _ Injector = new(singleInjector)

func NewSingleInjector() Injector {
	return &singleInjector{
		objs: make(map[interface{}]reflect.Value),
	}
}

type singleInjector struct {
	objs map[interface{}]reflect.Value
}

// Replace 替换 NewFunc
func (s *singleInjector) Replace(key interface{}, value interface{}) {
	s.objs[key] = reflect.ValueOf(value)
}

// Register 注册 NewFunc
func (s *singleInjector) Register(key interface{}, value interface{}) error {
	_, ok := s.objs[key]
	if ok {
		return errors.New(fmt.Sprintf("key duplicate: %v", key))
	}
	s.objs[key] = reflect.ValueOf(value)
	return nil
}

// Get 获取注册对象
func (s *singleInjector) Get(key string) (interface{}, error) {
	v, ok := s.objs[key]
	if ok {
		return v.Interface(), nil
	}
	return nil, ErrNotFound
}

// Remove 删除注册对象
func (s *singleInjector) Remove(key string) {
	delete(s.objs, key)
}

func (s *singleInjector) Clear() {
	s.objs = make(map[interface{}]reflect.Value)
}

func (s *singleInjector) InjectAll() {
	for _, v := range s.objs {
		s.Inject(v)
	}
}

func (s *singleInjector) Inject(v interface{}) {
	var value reflect.Value
	var ok bool
	if value, ok = v.(reflect.Value); !ok {
		value = reflect.ValueOf(v)
	}

loop:
	for {
		switch value.Kind() {
		case reflect.Ptr:
			value = value.Elem()
		case reflect.Interface:
			value = value.Elem()
		default:
			break loop
		}
	}

	if value.Kind() != reflect.Struct {
		return
	}
	s.injectValue(value)
}

func (s *singleInjector) injectValue(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		name := fieldType.Tag.Get(InjectorTag)
		if len(name) == 0 {
			if fieldType.Anonymous {
				s.injectValue(value.Field(i))
			}
			continue
		}
		temp, ok := s.objs[name]
		if ok {
			field := value.Field(i)
			if field.CanSet() {
				field.Set(temp)
			} else {
				field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				field.Set(temp)
			}
		}
	}
}
