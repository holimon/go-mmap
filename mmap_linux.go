package mmap

import (
	"errors"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

func Acquire[T DataAble](config Config) (m *Mmap[T], err error) {
	m = &Mmap[T]{}
	t := int(unsafe.Sizeof(T(0)))
	m.align = t / int(unsafe.Sizeof(uint8(0)))
	page := os.Getpagesize()
	if (config.Base%int64(page) != 0) && (!config.AlignPage) {
		return nil, errors.New("baseaddr must be a multiple of the system's page size")
	}
	base := config.Base / int64(page) * int64(page)
	offset := config.Base % int64(page)
	len := config.Length + int(offset)
	if len%t != 0 && config.AlignSize {
		len = len / t * t
	}
	if m.file, err = os.OpenFile(config.FileName, os.O_RDWR|os.O_SYNC, 0644); err != nil {
		return nil, err
	}
	if m.ref, err = syscall.Mmap(
		int(m.file.Fd()),
		base,
		len,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	); err != nil {
		m.file.Close()
		return nil, err
	}
	m.base = base
	m.length = len
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&m.ref))
	header.Len /= m.align
	header.Cap /= m.align
	m.data = *(*[]T)(unsafe.Pointer(&header))
	return m, nil
}

func (m *Mmap[T]) Release() {
	syscall.Munmap(m.ref)
	m.file.Close()
}
