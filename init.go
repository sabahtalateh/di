package di

import (
	"fmt"
	"runtime/debug"
)

func (c *Container) enterInit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return ErrInitialized
	}

	if c.initializing {
		return ErrInitialized
	}

	c.initializing = true

	return nil
}

func (c *Container) exitInit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initialized = true
	c.initializing = false
}

func InitE[T any](f func(*Container) (T, error)) withInitE[T] { return f }
func Init[T any](f func(*Container) T) withInit[T]            { return f }

func (c *Container) Init() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %s", toError(r), debug.Stack())
			if !recoverable(err) {
				panic(err)
			}
		}
	}()

	if err = c.enterInit(); err != nil {
		return err
	}

	for _, coord := range c.initOrder {
		comp, ok := c.components[coord]
		if !ok {
			continue
		}

		comp.val, err = comp.initFn(c)
		if err != nil {
			return err
		}
		comp.initFn = nil

		c.components[coord] = comp
	}

	c.exitInit()

	return err
}
