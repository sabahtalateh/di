package di

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStageType struct {
	val string
}

func Test_error_stage_set(t *testing.T) {
	c := NewContainer()

	err := Setup[testStageType](c,
		Init(func(c *Container) testStageType { return testStageType{} }),
		Stage("stage", func(ctx context.Context, s testStageType) error { return nil }),
		Stage("stage", func(ctx context.Context, s testStageType) error { return nil }),
	)

	require.ErrorIs(t, err, ErrStageSet)
}

func Test_context_in_all_stage_funcctions_cancelled_on_error_in_one(t *testing.T) {
	var (
		c      = NewContainer()
		aError = errors.New("error in A")
		bError error
		cError error
	)

	err := Setup[testStageType](c,
		Name("A"),
		Init(func(c *Container) testStageType { return testStageType{} }),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			return aError
		}),
	)
	require.NoError(t, err)

	err = Setup[testStageType](c,
		Name("B"),
		Init(func(c *Container) testStageType { return testStageType{} }),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			<-ctx.Done()
			bError = context.Cause(ctx)
			return bError
		}),
	)
	require.NoError(t, err)

	err = Setup[testStageType](c,
		Name("C"),
		Init(func(c *Container) testStageType { return testStageType{} }),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			<-ctx.Done()
			cError = context.Cause(ctx)
			return cError
		}),
	)
	require.NoError(t, err)

	err = c.Init()
	require.NoError(t, err)

	err = c.ExecStage(context.Background(), "stage")

	require.ErrorIs(t, err, aError)
	require.ErrorIs(t, err, ErrExecuteStage)

	require.ErrorIs(t, err, bError)
	require.ErrorIs(t, err, ErrExecuteStage)

	require.ErrorIs(t, err, cError)
	require.ErrorIs(t, err, ErrExecuteStage)
}

func Test_stage_function_receives_initialized_component(t *testing.T) {
	var (
		c    = NewContainer()
		aVal string
		bVal string
		cVal string
	)

	err := Setup[testStageType](c,
		Name("A"),
		Init(func(c *Container) testStageType {
			return testStageType{
				val: "A",
			}
		}),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			aVal = s.val
			return nil
		}),
	)
	require.NoError(t, err)

	err = Setup[testStageType](c,
		Name("B"),
		Init(func(c *Container) testStageType {
			return testStageType{
				val: "B",
			}
		}),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			bVal = s.val
			return nil
		}),
	)
	require.NoError(t, err)

	err = Setup[testStageType](c,
		Name("C"),
		Init(func(c *Container) testStageType {
			return testStageType{
				val: "C",
			}
		}),
		Stage("stage", func(ctx context.Context, s testStageType) error {
			cVal = s.val
			return nil
		}),
	)
	require.NoError(t, err)

	err = c.Init()
	require.NoError(t, err)

	err = c.ExecStage(context.Background(), "stage")
	require.NoError(t, err)

	require.Equal(t, "A", aVal)
	require.Equal(t, "B", bVal)
	require.Equal(t, "C", cVal)
}

func Test_error_exec_stage_before_init(t *testing.T) {
	c := NewContainer()
	err := c.ExecStage(context.Background(), "stage")
	require.ErrorIs(t, err, ErrNotInitialized)
}
