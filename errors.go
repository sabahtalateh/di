package di

import (
	"errors"
	"fmt"
)

var (
	ErrInitialized    = fmt.Errorf("initialized")
	ErrNotInitialized = fmt.Errorf("not initialized")
	ErrComponentSet   = fmt.Errorf("component set")
	ErrNameSet        = fmt.Errorf("name set")
	ErrInitSet        = fmt.Errorf("init function set")
	ErrInitNotSet     = fmt.Errorf("init function not set")
	ErrStageSet       = fmt.Errorf("stage set")
	ErrExecuteStage   = fmt.Errorf("execute stage")
	ErrNotFound       = fmt.Errorf("not found")

	recoverableErrs = []error{
		ErrInitialized,
		ErrNotInitialized,
		ErrComponentSet,
		ErrNameSet,
		ErrInitSet,
		ErrInitNotSet,
		ErrStageSet,
		ErrExecuteStage,
		ErrNotFound,
	}
)

func recoverable(err error) bool {
	for _, e := range recoverableErrs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func toError(x any) (err error) {
	switch x := x.(type) {
	case error:
		err = x
	default:
		err = fmt.Errorf("%s", x)
	}
	return err
}
