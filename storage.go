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
	storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error
	clearCheckpoint(taskId string, msgId string, cpId string) error
	getCheckpointMessages(taskId string, cpId string) ([][]byte, error)

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

func storeCheckpoint[T any](s Storage, taskId string, cpId string, msg *Message[T]) error {
	b, err := msg.messageToBytes()
	if err != nil {
		return err
	}

	s.storeCheckpoint(taskId, msg.Id, cpId, b)
	return nil
}

func recoverCheckpoint[T any](s Storage, taskId string, cpId string, onMessage func(m *Message[T])) error {
	msgs, err := s.getCheckpointMessages(taskId, cpId)
	if err != nil {
		return err
	}

	for i := range msgs {
		msg, err := bytesToMessage[T](&msgs[i])
		if err != nil {
			return err
		}

		onMessage(msg)
	}

	return nil
}

func startNewEmptyWindow(s Storage, cpId string, key string, start ...time.Time) (windowMeta, error) {
	return s.startNewEmptyWindow(cpId, key, start...)
}

func startNewWindow[T any](s Storage, cpId string, key string, msg Message[T], start ...time.Time) (windowMeta, error) {
	return s.startNewWindow(cpId, key, msg, start...)
}

func pushMessageToWindow[T any](s Storage, cpId string, key string, winId string, msg Message[T]) (int, error) {
	return s.pushMessageToWindow(cpId, key, winId, msg)
}

func getWindowMetadata(s Storage, cpId string, key string, winId string) (windowMeta, error) {
	return s.getWindowMetadata(cpId, key, winId)
}

func getWindowsMetadataByKey(s Storage, cpId string, key string) ([]windowMeta, error) {
	return s.getWindowsMetadataByKey(cpId, key)
}

func closeWindow[T any](s Storage, cpId string, key string, winId string, watermark time.Duration, onFlush func(items []Message[T])) error {
	err := s.closeWindow(cpId, key, winId)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(watermark)
		items, err := flushWindow[T](s, cpId, key, winId)
		if err != nil {
			panic(err)
		}

		onFlush(items)
	}()

	return nil
}

func flushWindow[T any](s Storage, cpId string, key string, winId string) ([]Message[T], error) {
	m, err := s.flushWindow(cpId, key, winId)
	if err != nil {
		return nil, err
	}

	msgs := make([]Message[T], len(m))

	for i := range m {
		msgs[i] = m[i].(Message[T])
	}

	return msgs, nil
}

func destroyWindow(s Storage, cpId string, key string, winId string) error {
	return s.destroyWindow(cpId, key, winId)
}

func getWindowSize(s Storage, cpId string, key string, winId string) (int, error) {
	return s.getWindowSize(cpId, key, winId)
}

func getKeys(s Storage, cpId string) ([]string, error) {
	return s.getKeys(cpId)
}
