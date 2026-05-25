package httpx

import (
	"fmt"
	"sync"
)

var (
	config configuration
	rwm    *sync.RWMutex = &sync.RWMutex{}
)

type configuration struct {
	maximumRequestBodySize          int64
	problemForInvalidJSON           func(error) *Problem
	problemForLengthRequired        func() *Problem
	problemForRequestEntityTooLarge func() *Problem
	problemForUnexpectedError       func(error) *Problem
	problemForUnprocessableEntity   func(error) *Problem
	problemForUnsupportedMediaType  func(string, string) *Problem
}

func UseMaximumRequestBodySize(size int64) {
	rwm.Lock()
	defer rwm.Unlock()

	config.maximumRequestBodySize = size
}

func UseProblemConstructorForInvalidJSON(fn func(error) *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForInvalidJSON = fn
}

func UseProblemConstructorForLengthRequired(fn func() *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForLengthRequired = fn
}

func UseProblemConstructorForRequestEntityTooLarge(fn func() *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForRequestEntityTooLarge = fn
}

func UseProblemConstructorForUnexpectedError(fn func(error) *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForUnexpectedError = fn
}

func UseProblemConstructorForUnprocessableEntity(fn func(error) *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForUnprocessableEntity = fn
}

func UseProblemConstructorForUnsupportedMediaType(fn func(string, string) *Problem) {
	rwm.Lock()
	defer rwm.Unlock()

	config.problemForUnsupportedMediaType = fn
}

func (c configuration) MaximumRequestBodySize() int64 {
	rwm.RLock()
	defer rwm.RUnlock()

	return c.maximumRequestBodySize
}

func (c configuration) ProblemForInvalidJSON(err error) *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForInvalidJSON == nil {
		return &Problem{
			Type:   "InvalidJSON",
			Detail: "The request body contains invalid JSON.",
		}
	}

	return c.problemForInvalidJSON(err)
}

func (c configuration) ProblemForLengthRequired() *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForLengthRequired == nil {
		return &Problem{
			Type:   "LengthRequired",
			Detail: "A request body is required but no Content-Length was provided.",
		}
	}

	return c.problemForLengthRequired()
}

func (c configuration) ProblemForRequestEntityTooLarge() *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForRequestEntityTooLarge == nil {
		return &Problem{
			Type:   "RequestEntityTooLarge",
			Detail: fmt.Sprintf("The request body is too large. Maximum permitted size is %d bytes.", c.maximumRequestBodySize),
		}
	}

	return c.problemForRequestEntityTooLarge()
}

func (c configuration) ProblemForUnexpectedError(err error) *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForUnexpectedError == nil {
		return &Problem{
			Type:   "UnexpectedError",
			Detail: "An unexpected error prevented the operation from completing.",
		}
	}

	return c.problemForUnexpectedError(err)
}

func (c configuration) ProblemForUnprocessableEntity(err error) *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForUnprocessableEntity == nil {
		return &Problem{
			Type:   "UnprocessableEntity",
			Detail: err.Error(),
		}
	}

	return c.problemForUnprocessableEntity(err)
}

func (c configuration) ProblemForUnsupportedMediaType(expected string, actual string) *Problem {
	rwm.RLock()
	defer rwm.RUnlock()

	if c.problemForUnsupportedMediaType == nil {
		return &Problem{
			Type:   "UnsupportedMediaType",
			Detail: fmt.Sprintf("Expected Content-Type '%s' but received '%s'.", expected, actual),
		}
	}

	return c.problemForUnsupportedMediaType(expected, actual)
}
