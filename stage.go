package di

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
)

func stageFn[T any](c *Container, coord coordinate, fn func(context.Context, T) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		comp, ok := c.components[coord]
		if !ok {
			// no way to delete component. impossible case
			return nil
		}

		t, ok := comp.val.(T)
		if !ok {
			// no way have type other then T for coord. impossible case
			return nil
		}

		return fn(ctx, t)
	}
}

func (c *Container) checkExecStage() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return ErrNotInitialized
	}

	if c.initializing {
		return ErrNotInitialized
	}

	return nil
}

func Stage[T any](name string, fn func(context.Context, T) error) withStage[T] {
	return withStage[T]{name: name, fn: fn}
}

func (c *Container) ExecStage(ctx context.Context, name string) error {
	if err := c.checkExecStage(); err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	ctx, cnl := context.WithCancelCause(ctx)
	defer cnl(nil)
	for _, stageFn := range c.stages[name] {
		fn := stageFn
		eg.Go(func() (err error) {
			if err = fn(ctx); err != nil {
				err = fmt.Errorf("%w: %s: %w: %s", ErrExecuteStage, name, err, debug.Stack())
				cnl(err)
			}
			return err
		})
	}

	return eg.Wait()
}
