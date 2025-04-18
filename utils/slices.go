package utils

import "github.com/pkg/errors"

func RemoveIndex[T any](s []T, index int) ([]T, error) {
	if len(s) <= index {
		return make([]T, 0), errors.New("未找到此内容")
	}
	if len(s) == index-1 {
		return s[:index], nil
	}
	return append(s[:index], s[index+1:]...), nil
}
