package stringsvc

import (
	"context"
	"errors"
)

// Service provides operations on strings.
type Service interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

func New() *StringService {
	return &StringService{}
}

// StringService is a concrete implementation of Service
type StringService struct{}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
