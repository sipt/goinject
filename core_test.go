package goinject

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInject(t *testing.T) {
	a := &A{}
	b := &B{}
	injector := DefaultInjector()
	Register(func(injector Injector) error {
		err := injector.Register("value.a", a)
		assert.NoError(t, err)
		err = injector.Register("value.a", a)
		assert.Error(t, ErrKeyDuplicate)
		err = injector.Register("string.a", "string.a")
		assert.NoError(t, err)
		err = injector.Register("value.b", b)
		assert.NoError(t, err)
		err = injector.Register("value.b.GetName", b.GetName)
		assert.NoError(t, err)
		return nil
	})
	err := InjectAll(injector)
	assert.NoError(t, err)
	assert.Equal(t, "string.a", "string.a")
	assert.NoError(t, err)
	assert.Equal(t, "b", a.ib.GetName())
	assert.NoError(t, err)
	assert.Equal(t, "b", a.b.GetName())
	assert.NoError(t, err)
	assert.Equal(t, "b", a.bf())
	assert.NoError(t, err)

	i, err := injector.Get("value.b")
	assert.NoError(t, err)
	ib, ok := i.(IB)
	assert.True(t, ok)
	assert.Equal(t, "b", ib.GetName())

	ib = &B{}
	injector.Inject(ib)
	assert.Equal(t, "string.a", ib.(*B).a)
}

type A struct {
	a string `inject:"string.a"`
	b *B `inject:"value.b"`
	ib IB `inject:"value.b"`
	bf func()string `inject:"value.b.GetName"`
}

type B struct {
	a string `inject:"string.a"`
}

func (*B)GetName() string {
	return "b"
}

type IB interface {
	GetName() string
}
