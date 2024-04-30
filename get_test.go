package di

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type getTestType struct{}

func Test_GetE(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() (*Container, error)
		getOpts    []getOpt[getTestType]
		want       getTestType
		wantGetErr error
	}{
		{
			name:  "error name set",
			setup: func() (*Container, error) { return NewContainer(), nil },
			getOpts: []getOpt[getTestType]{
				Name("A"),
				Name("B"),
			},
			wantGetErr: ErrNameSet,
		},
		{
			name:  "error not found when container empty",
			setup: func() (*Container, error) { return NewContainer(), nil },
			getOpts: []getOpt[getTestType]{
				Name("A"),
			},
			wantGetErr: ErrNotFound,
		},
		{
			name: "error not found when names not match",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[getTestType](c,
					Name("A"),
					Init(func(c *Container) getTestType {
						return getTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[getTestType]{
				Name("B"),
			},
			wantGetErr: ErrNotFound,
		},
		{
			name: "get named",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[getTestType](c,
					Name("A"),
					Init(func(c *Container) getTestType {
						return getTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[getTestType]{
				Name("A"),
			},
			want: getTestType{},
		},
		{
			name: "get unnamed",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[getTestType](c,
					Init(func(c *Container) getTestType {
						return getTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[getTestType]{},
			want:    getTestType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := tt.setup()
			require.NoError(t, err)

			err = c.Init()
			require.NoError(t, err)

			got, err := GetE(c, tt.getOpts...)
			if tt.wantGetErr != nil {
				require.ErrorIs(t, err, tt.wantGetErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_GetE_pointer_type(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() (*Container, error)
		getOpts    []getOpt[*getTestType]
		want       *getTestType
		wantGetErr error
	}{
		{
			name: "get named",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*getTestType](c,
					Name("A"),
					Init(func(c *Container) *getTestType {
						return &getTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[*getTestType]{
				Name("A"),
			},
			want: &getTestType{},
		},
		{
			name: "get unnamed",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*getTestType](c,
					Init(func(c *Container) *getTestType {
						return &getTestType{}
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[*getTestType]{},
			want:    &getTestType{},
		},
		{
			name: "get nil",
			setup: func() (*Container, error) {
				c := NewContainer()
				err := Setup[*getTestType](c,
					Init(func(c *Container) *getTestType {
						return nil
					}),
				)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			getOpts: []getOpt[*getTestType]{},
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := tt.setup()
			require.NoError(t, err)

			err = c.Init()
			require.NoError(t, err)

			got, err := GetE(c, tt.getOpts...)
			if tt.wantGetErr != nil {
				require.ErrorIs(t, err, tt.wantGetErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_get_before_init(t *testing.T) {
	c := NewContainer()

	_, err := GetE[getTestType](c)
	require.ErrorIs(t, err, ErrNotInitialized)
}

func Test_get_err_hint_1(t *testing.T) {
	c := NewContainer()

	// setup pointer
	Setup[*getTestType](c,
		Init(func(c *Container) *getTestType {
			return &getTestType{}
		}),
	)
	err := c.Init()
	require.NoError(t, err)

	// get not pointer
	_, err = GetE[getTestType](c)
	require.Error(t, err, ErrNotFound)
	require.True(t, strings.Contains(err.Error(), "found component (*di.getTestType, (Unnamed))"))
}

func Test_get_err_hint_2(t *testing.T) {
	c := NewContainer()

	// setup not pointer
	Setup[getTestType](c,
		Init(func(c *Container) getTestType {
			return getTestType{}
		}),
	)
	err := c.Init()
	require.NoError(t, err)

	// get pointer
	_, err = GetE[*getTestType](c)
	require.Error(t, err, ErrNotFound)
	require.True(t, strings.Contains(err.Error(), "found component (di.getTestType, (Unnamed))"))
}

func Test_get_stack_trace(t *testing.T) {
	type getTestType2 struct {
		gtt *getTestType
	}

	c := NewContainer()
	err := Setup[*getTestType2](c,
		Init(func(c *Container) *getTestType2 {
			return &getTestType2{
				gtt: Get[*getTestType](c),
			}
		}),
	)
	require.NoError(t, err)

	err = c.Init()
	// test err contains stack trace with file path
	require.Contains(t, err.Error(), "github.com/sabahtalateh/di/get_test.go:")
}
