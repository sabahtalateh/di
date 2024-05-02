package di

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type containerOpt interface {
	applyContainerOpt(*Container)
}

type coordinate struct {
	type_ reflect.Type
	name  string
}

func (c coordinate) String() string {
	name := c.name
	if name == "" {
		name = "(Unnamed)"
	}
	return fmt.Sprintf("(%s, %s)", c.type_, name)
}

type component struct {
	// initFn also used to indicate if component initialized
	// if initFn is not nil component not initialized yet
	// if initFn is nil component initialized
	initFn func(*Container) (any, error)
	val    any
}

type Container struct {
	mu           sync.Mutex
	initializing bool
	initialized  bool

	initOrder  []coordinate
	components map[coordinate]component
	stages     map[string][]func(context.Context) error
}

func NewContainer(opts ...containerOpt) *Container {
	c := &Container{
		components: make(map[coordinate]component),
		stages:     make(map[string][]func(context.Context) error),
	}

	for _, o := range opts {
		o.applyContainerOpt(c)
	}

	return c
}
