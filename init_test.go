package di

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type initTestType struct{}

type initTestType2 struct {
	itt initTestType
}

func Test_Init(t *testing.T) {
	tests := []struct {
		name            string
		setup           func() (*Container, error)
		wantInitErr     error
		wantErrContains []string
	}{
		{
			name: "error get component component which is not set",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[initTestType2](c,
					Init(func(c *Container) initTestType2 {
						return initTestType2{
							itt: Get[initTestType](c),
						}
					}),
				)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantInitErr: ErrNotFound,
		},
		{
			name: "error disordered setup",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[initTestType2](c,
					Init(func(c *Container) initTestType2 {
						return initTestType2{
							itt: Get[initTestType](c),
						}
					}),
				)
				if err != nil {
					return nil, err
				}

				err = Setup[initTestType](c,
					Init(func(c *Container) initTestType {
						return initTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantInitErr: ErrDisordered,
			wantErrContains: []string{
				"(di.initTestType, (Unnamed)) must be set before parent component",
			},
		},
		{
			name: "error init after init",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := c.Init()
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantInitErr: ErrInitialized,
		},
		{
			name: "ok",
			setup: func() (*Container, error) {
				c := NewContainer()

				err := Setup[initTestType](c,
					Init(func(c *Container) initTestType {
						return initTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}

				err = Setup[initTestType2](c,
					Init(func(c *Container) initTestType2 {
						return initTestType2{
							itt: Get[initTestType](c),
						}
					}),
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
			c, err := tt.setup()
			require.NoError(t, err)

			err = c.Init()
			if tt.wantInitErr != nil {
				require.ErrorIs(t, err, tt.wantInitErr)
				for _, str := range tt.wantErrContains {
					require.Contains(t, err.Error(), str)
				}
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_user_panic_in_init_not_recovered(t *testing.T) {
	c := NewContainer()
	err := Setup[initTestType](c,
		Init(func(c *Container) initTestType {
			panic("some panic")
		}),
	)
	require.NoError(t, err)

	var panic string
	func() {
		defer func() {
			if r := recover(); r != nil {
				panic = fmt.Sprintf("%s", r)
			}
		}()
		_ = c.Init()
	}()

	require.True(t, strings.HasPrefix(panic, "some panic"))
}
