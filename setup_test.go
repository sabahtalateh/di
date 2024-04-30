package di

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type setupTestType struct{}
type setupTestType2 struct{}

func Test_Setup(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (*Container, error)
		wantSetupErr error
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
			wantSetupErr: ErrInitNotSet,
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
			wantSetupErr: ErrInitSet,
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
			wantSetupErr: ErrInitSet,
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
			wantSetupErr: ErrInitSet,
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
			wantSetupErr: ErrInitSet,
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
			wantSetupErr: ErrNameSet,
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
			wantSetupErr: ErrNameSet,
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
			wantSetupErr: ErrNameSet,
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
			wantSetupErr: ErrNameSet,
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
			wantSetupErr: ErrNameSet,
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
			wantSetupErr: ErrComponentSet,
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
			wantSetupErr: ErrInitNotSet,
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
			wantSetupErr: ErrInitialized,
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
				return
			}
			require.NoError(t, err)
		})
	}
}
