package di

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
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
		return fmt.Errorf("%w: %s", ErrInitialized, debug.Stack())
	}

	if c.initializing {
		return fmt.Errorf("%w: %s", ErrInitialized, debug.Stack())
	}

	return nil
}

func processSetupOpts[T any](opts ...setupOpt[T]) (
	// name
	string,
	// init function
	func(*Container) (any, error),
	// stages functions
	map[string]func(context.Context, T) error,
	error,
) {
	var (
		t        T
		name     = ""
		nameSet  = false
		initFn   func(*Container) (any, error)
		stageFns = make(map[string]func(context.Context, T) error)
	)

	for _, o := range opts {
		switch o := o.(type) {
		case withName:
			if nameSet {
				return "", nil, nil, fmt.Errorf("%w: %s", ErrNameSet, debug.Stack())
			}

			name = string(o)
			nameSet = true
		case withInitE[T]:
			if o == nil {
				return "", nil, nil, fmt.Errorf("%w: for type (%s): %s", ErrInitNotSet, reflect.TypeOf(&t).Elem(), debug.Stack())
			}
			if initFn != nil {
				return "", nil, nil, fmt.Errorf("%w: for type (%s): %s", ErrInitSet, reflect.TypeOf(&t).Elem(), debug.Stack())
			}

			initFn = func(c *Container) (any, error) { return o(c) }
		case withInit[T]:
			if o == nil {
				return "", nil, nil, fmt.Errorf("%w: for type (%s): %s", ErrInitNotSet, reflect.TypeOf(&t).Elem(), debug.Stack())
			}
			if initFn != nil {
				return "", nil, nil, fmt.Errorf("%w: for type (%s): %s", ErrInitSet, reflect.TypeOf(&t).Elem(), debug.Stack())
			}

			initFn = func(c *Container) (any, error) { return o(c), nil }
		case withStage[T]:
			if _, ok := stageFns[o.name]; ok {
				return "", nil, nil, fmt.Errorf("%w: %s", ErrStageSet, debug.Stack())
			}
			if o.fn == nil {
				return "", nil, nil, fmt.Errorf("%w: %s", ErrStageNotSet, debug.Stack())
			}

			stageFns[o.name] = o.fn
		}
	}

	if initFn == nil {
		return "", nil, nil, fmt.Errorf("%w: for type (%s): %s", ErrInitNotSet, reflect.TypeOf(&t).Elem(), debug.Stack())
	}

	return name, initFn, stageFns, nil
}

func Setup[T any](c *Container, opts ...setupOpt[T]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.checkSetup(); err != nil {
		return err
	}

	name, initFn, stageFns, err := processSetupOpts(opts...)
	if err != nil {
		return err
	}

	var (
		t T
	)

	coord := coordinate{
		type_: reflect.TypeOf(&t).Elem(),
		name:  name,
	}

	if _, ok := c.components[coord]; ok {
		return fmt.Errorf("%w: %s", ErrComponentSet, debug.Stack())
	}

	c.components[coord] = component{
		initFn: initFn, // set to nil after initialization
	}

	c.initOrder = append(c.initOrder, coord)

	for name, fn := range stageFns {
		c.stages[name] = append(c.stages[name], stageFn(c, coord, fn))
	}

	return nil
}
