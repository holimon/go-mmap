package mmap

import (
	"errors"
	"os"
)

type DataAble interface {
	uint8 | uint16 | uint32 | uint64
}

type Mmap[T DataAble] struct {
	file   *os.File
	ref    []byte
	base   int64
	length int
	data   []T
	align  int
}

type Config struct {
	AlignPage bool
	AlignSize bool
	Base      int64
	Length    int
	FileName  string
}

func (m *Mmap[T]) verify(addr int64, bits ...int) error {
	offset := int(addr - m.base)
	if offset%m.align != 0 || offset >= m.length || (offset < 0) {
		return errors.New("invalid memory address")
	}
	for _, b := range bits {
		if b >= m.align*8 || b < 0 {
			return errors.New("invalid bit")
		}
	}
	return nil
}

func (m *Mmap[T]) BaseAddress() int64 {
	return m.base
}

func (m *Mmap[T]) TotalLength() int {
	return m.length
}

func (m *Mmap[T]) MemoryRead(addr int64) (T, error) {
	if err := m.verify(addr); err != nil {
		return 0, err
	}
	return m.data[int(addr-m.base)/m.align], nil
}

func (m *Mmap[T]) MemoryWrite(addr int64, val T) error {
	if err := m.verify(addr); err != nil {
		return err
	}
	m.data[int(addr-m.base)/m.align] = val
	return nil
}

func (m *Mmap[T]) MemorySpecialMask(addr int64, bits ...int) error {
	if err := m.verify(addr, bits...); err != nil {
		return err
	}
	for _, b := range bits {
		m.data[int(addr-m.base)/m.align] |= T(1 << b)
	}
	return nil
}
func (m *Mmap[T]) MemorySpecialClear(addr int64, bits ...int) error {
	if err := m.verify(addr, bits...); err != nil {
		return err
	}
	for _, b := range bits {
		m.data[int(addr-m.base)/m.align] &= ^T(1 << b)
	}
	return nil
}
func (m *Mmap[T]) MemorySpecialNegate(addr int64, bits ...int) error {
	if err := m.verify(addr, bits...); err != nil {
		return err
	}
	for _, b := range bits {
		m.data[int(addr-m.base)/m.align] ^= T(1 << b)
	}
	return nil
}
