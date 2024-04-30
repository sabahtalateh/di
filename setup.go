package di

import (
	"context"
	"fmt"
	"reflect"
)

type setupOpt[T any] interface {
	setupOpt()
}

type withInitE[T any] func(*Container) (T, error)
type withInit[T any] func(*Container) T
type withStage[T any] struct {
	name string
	fn   func(context.Context, T) error
}

func (o withName) setupOpt()   {}
func (withInitE[T]) setupOpt() {}
func (withInit[T]) setupOpt()  {}
func (withStage[T]) setupOpt() {}

func (c *Container) checkSetup() error {
	if c.initialized {
		return ErrInitialized
	}

	if c.initializing {
		return ErrInitialized
	}

	return nil
}

func Setup[T any](c *Container, opts ...setupOpt[T]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.checkSetup(); err != nil {
		return err
	}

	var (
		t        T
		name     = ""
		nameSet  = false
		initFn   func(c *Container) (any, error)
		stageFns = make(map[string]func(context.Context, T) error)
	)

	// process options
	{
		for _, o := range opts {
			switch o := o.(type) {
			case withName:
				if nameSet {
					return ErrNameSet
				}
				name = string(o)
				nameSet = true
			case withInitE[T]:
				if initFn != nil {
					return ErrInitSet
				}
				initFn = func(c *Container) (any, error) { return o(c) }
			case withInit[T]:
				if initFn != nil {
					return ErrInitSet
				}
				initFn = func(c *Container) (any, error) { return o(c), nil }
			case withStage[T]:
				if _, ok := stageFns[o.name]; ok {
					return fmt.Errorf("%w: %s", ErrStageSet, o.name)
				}
				stageFns[o.name] = o.fn
			}
		}

		if initFn == nil {
			return fmt.Errorf("%w: for type (%s)", ErrInitNotSet, reflect.TypeOf(&t).Elem())
		}
	}

	// setup component
	{
		coord := coordinate{
			type_: reflect.TypeOf(&t).Elem(),
			name:  name,
		}

		if _, ok := c.components[coord]; ok {
			if coord.name == "" {
				coord.name = "(Unnamed)"
			}
			return fmt.Errorf("%w (%s, %s)", ErrComponentSet, coord.type_, coord.name)
		}

		c.components[coord] = component{
			initFn: initFn, // set to nil after initialization
			val:    nil,    // set to init function return value after initialization
		}

		for name, fn := range stageFns {
			c.stages[name] = append(c.stages[name], stageFn(c, coord, fn))
		}
	}

	return nil
}
