package models

type Clone[T any] interface {
	Clone() T
}

type GetValue func(string) string
