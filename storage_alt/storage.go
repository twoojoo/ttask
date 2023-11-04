package storage

type Storage[T any] interface {
	StartNewWindow(key string, elem T, start ...int64) WindowMeta
	StartNewEmptyWindow(key string, start ...int64) WindowMeta
	PushItemToWindow(k string, id string, item T) int
	GetWindowsMetadata(k string) []WindowMeta
	GetWindowMetadata(k string, id string) WindowMeta
	CloseWindow(k string, id string)
	FlushWindow(k string, id string) []T
	DestroyWindow(k string, id string)
	GetWindowSize(k string, id string) int
	GetKeys() []string

	StoreCheckpoint(mId string, cId string, msg T)
	ClearCheckpoint(mId string)
	GetCheckpointMessages(cId string) []T
}
