package utils

import "log"

type TODOStruct struct{}

func (TODOStruct) String() string {
	panic("not implemented")
}

func Debug[T any](msg string, v T) T {
	log.Printf("%s: %v", msg, v)
	return v
}
