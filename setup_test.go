package di

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type setupTestType struct{}
type setupTestType2 struct{}

func Test_Setup(t *testing.T) {
	tests := []struct {
		name            string
		setup           func() (*Container, error)
		wantSetupErr    error
		wantErrContains []string
	}{
		{
			name: "error no init",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[setupTestType](c)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitNotSet,
			wantErrContains: []string{"TEST"},
		},
		{
			name: "error init function is nil",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[initTestType](c,
					InitE[initTestType](nil),
				)

				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantSetupErr:    ErrInitNotSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error init set 1",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error init set 2",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
					InitE(func(c *Container) (*setupTestType, error) { return new(setupTestType), nil }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error init set 3",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					InitE(func(c *Container) (*setupTestType, error) { return new(setupTestType), nil }),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error init set 4",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					InitE(func(c *Container) (*setupTestType, error) { return new(setupTestType), nil }),
					InitE(func(c *Container) (*setupTestType, error) { return new(setupTestType), nil }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error name set 1",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name("A"),
					Name("B"),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrNameSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error name set 2",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name(""),
					Name("A"),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrNameSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error name set 3",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name("A"),
					Name(""),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrNameSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error name set 4",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name(""),
					Name(""),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrNameSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error same name set",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name("A"),
					Name("A"),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrNameSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error component set",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name("A"),
					Init(func(c *Container) *setupTestType {
						return &setupTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}

				err = Setup[*setupTestType](c,
					Name("A"),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrComponentSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error stage function not set",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
					Stage[*setupTestType]("stage", nil),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrStageNotSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error type mismatch",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Init(func(c *Container) *setupTestType2 { return new(setupTestType2) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			wantSetupErr:    ErrInitNotSet,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "error setup after init",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := c.Init()
				if err != nil {
					return nil, err
				}

				err = Setup[setupTestType](c,
					Init(func(c *Container) setupTestType { return setupTestType{} }),
				)

				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantSetupErr:    ErrInitialized,
			wantErrContains: []string{"di/setup_test.go"},
		},
		{
			name: "ok 1",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
		{
			name: "ok 2",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*setupTestType](c,
					Name("A"),
					Init(func(c *Container) *setupTestType { return new(setupTestType) }),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.setup()
			if tt.wantSetupErr != nil {
				require.ErrorIs(t, err, tt.wantSetupErr)
				for _, str := range tt.wantErrContains {
					require.Contains(t, err.Error(), str)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
