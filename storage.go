package ttask

import (
	"time"
)

type window struct {
	id    string
	start time.Time
	end   time.Time
	last  time.Time
	elems []any
}

type windowMeta struct {
	Id    string
	Start time.Time
	End   time.Time
	Last  time.Time
}

type Storage interface {
	//checkpointing
	storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error
	clearCheckpoint(taskId string, msgId string, cpId string) error
	getCheckpointMessages(taskId string, cpId string) ([][]byte, error)

	//windowing
	startNewEmptyWindow(cpId string, key string, start ...time.Time) (windowMeta, error)
	startNewWindow(cpId string, key string, msg any, start ...time.Time) (windowMeta, error)
	pushMessageToWindow(cpId string, key string, winId string, msg any) (int, error)
	getWindowMetadata(cpId string, key string, winId string) (windowMeta, error)
	getWindowsMetadataByKey(cpId string, key string) ([]windowMeta, error)
	closeWindow(cpId string, key string, winId string) error
	flushWindow(cpId string, key string, winId string) ([]any, error)
	destroyWindow(cpId string, key string, winId string) error
	getWindowSize(cpId string, key string, winId string) (int, error)
	getKeys(cpId string) ([]string, error)
}
