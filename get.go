package di

import (
	"errors"
	"fmt"
	"reflect"
)

type getOpt[T any] interface {
	getOpt()
}

func (o withName) getOpt() {}

func (c *Container) checkGet() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized && !c.initializing {
		return ErrNotInitialized
	}

	return nil
}

func Get[T any](c *Container, opts ...getOpt[T]) T {
	t, err := GetE(c, opts...)
	if err != nil {
		// di errors converted into error when happens while init. see Init
		panic(err)
	}
	return t
}

func GetE[T any](c *Container, opts ...getOpt[T]) (T, error) {
	var (
		t T
	)

	if err := c.checkGet(); err != nil {
		return t, err
	}

	var (
		nameSet = false
		name    = ""
	)

	for _, o := range opts {
		switch o := o.(type) {
		case withName:
			if nameSet {
				return t, ErrNameSet
			}
			name = string(o)
			nameSet = true
		}
	}

	coord := coordinate{
		type_: reflect.TypeOf(&t).Elem(),
		name:  name,
	}

	val, err := c.initComponent(coord)
	if errors.Is(err, ErrNotFound) {
		return t, notFoundHint[T](c, coord.name)
	}
	if err != nil {
		return t, err
	}

	t, ok := val.(T)
	if !ok {
		return t, ErrNotFound
	}

	return t, nil
}

func notFoundHint[T any](c *Container, name string) error {
	tryCoord := coordinate{name: name}

	var t T
	pt := &t
	type_ := reflect.TypeOf(pt).Elem()

	if type_.Kind() == reflect.Pointer {
		tryCoord.type_ = type_.Elem()
	} else {
		tryCoord.type_ = reflect.TypeOf(&pt).Elem()
	}

	if _, ok := c.components[tryCoord]; ok {
		name := tryCoord.name
		if name == "" {
			name = "(Unnamed)"
		}
		return fmt.Errorf("%w: found component (%s, %s)", ErrNotFound, tryCoord.type_, name)
	}

	return ErrNotFound
}
