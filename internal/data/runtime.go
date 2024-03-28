package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d min", r)
	quotedValue := strconv.Quote(jsonValue)
	return []byte(quotedValue), nil
}

func (r *Runtime) UnmarshalJSON(json []byte) error {
	uqJson, err := strconv.Unquote(string(json))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(uqJson, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)
	return nil
}
