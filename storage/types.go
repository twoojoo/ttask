package storage

type window[T any] struct {
	id       string
	metadata map[string]int64
	elems    []T
}

type WindowMeta struct {
	Id       string
	Metadata map[string]int64
}