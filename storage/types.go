package storage

type window[T any] struct {
	id       string
	metadata map[string]int64
	start    int64
	end      int64
	elems    []T
}

type WindowMeta struct {
	Id       string
	Metadata map[string]int64
	Start    int64
	End      int64
}
