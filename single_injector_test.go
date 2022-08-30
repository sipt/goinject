package goinject

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnonymousInject(t *testing.T) {
	c := &C{}
	injector := DefaultInjector()
	Register(func(injector Injector) error {
		err := injector.Register("string.d", "string.d")
		assert.NoError(t, err)
		return nil
	})
	err := InjectAll(injector)
	assert.NoError(t, err)
	injector.Inject(c)
	assert.Equal(t, c.d, "string.d")
}

func TestProtectedInject(t *testing.T) {
	c := &C{}
	injector := DefaultInjector()
	Register(func(injector Injector) error {
		err := injector.Register("string.d", "string.d")
		assert.NoError(t, err)
		return nil
	})
	err := injector.Register("@protected_var", c)
	assert.NoError(t, err)
	err = InjectAll(injector)
	assert.NoError(t, err)

	assert.Equal(t, "", c.d)
}

type C struct {
	D
}

type D struct {
	d string `inject:"string.d"`
}
