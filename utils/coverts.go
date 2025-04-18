package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

func ParseTtl(ttl string) (uint32, error) {
	t, err := strconv.Atoi(ttl)
	if err != nil {
		return 0, err
	}
	if t < 0 {
		return 0, errors.New("invalid ttl")
	}
	return uint32(t), nil
}

func AtoUint32(src string) (uint32, error) {
	r, err := strconv.Atoi(src)
	if err != nil {
		return 0, err
	}
	if r < 0 {
		return 0, fmt.Errorf("%s is not a valid uint32", src)
	}
	return uint32(r), nil
}
