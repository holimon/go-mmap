package mmap

import (
	"errors"
)

var ErrorSupport = errors.New("not support")

func Acquire[T DataAble](config Config) (m *Mmap[T], err error) {
	return nil, ErrorSupport
}

func (m *Mmap[T]) Release() {
}
