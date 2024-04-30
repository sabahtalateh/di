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
		name        string
		setup       func() (*Container, error)
		wantInitErr error
	}{
		{
			name: "error get component while initializing",
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
			name: "error call setup from init",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[initTestType](c,
					InitE(func(c *Container) (initTestType, error) {
						if err := Setup[initTestType](c); err != nil {
							return initTestType{}, err
						}
						return initTestType{}, nil
					}),
				)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			wantInitErr: ErrInitialized,
		},
		{
			name: "init after init",
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
