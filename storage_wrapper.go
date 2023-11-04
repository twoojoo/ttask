package ttask

import (
	"time"
)

type storageWrapper[T any] struct {
	taskInner *Inner
	storage   Storage
}

func newStorageWrapper[T any](i *Inner) *storageWrapper[T] {
	return &storageWrapper[T]{
		taskInner: i,
		storage:   i.storage,
	}
}

func (sw *storageWrapper[T]) storeCheckpoint(taskId string, cpId string, msg *Message[T]) error {
	b, err := msg.messageToBytes()
	if err != nil {
		return err
	}

	sw.storage.storeCheckpoint(taskId, msg.Id, cpId, b)
	return nil
}

func (sw *storageWrapper[T]) recoverCheckpoint(taskId string, cpId string, onMessage func(m *Message[T])) error {
	msgs, err := sw.storage.getCheckpointMessages(taskId, cpId)
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

func (sw *storageWrapper[T]) startNewEmptyWindow(cpId string, key string, start ...time.Time) (windowMeta, error) {
	return sw.storage.startNewEmptyWindow(cpId, key, start...)
}

func (sw *storageWrapper[T]) startNewWindow(cpId string, key string, msg Message[T], start ...time.Time) (windowMeta, error) {
	return sw.storage.startNewWindow(cpId, key, msg, start...)
}

func (sw *storageWrapper[T]) pushMessageToWindow(cpId string, key string, winId string, msg Message[T]) (int, error) {
	return sw.storage.pushMessageToWindow(cpId, key, winId, msg)
}

func (sw *storageWrapper[T]) getWindowMetadata(cpId string, key string, winId string) (windowMeta, error) {
	return sw.storage.getWindowMetadata(cpId, key, winId)
}

func (sw *storageWrapper[T]) getWindowsMetadataByKey(cpId string, key string) ([]windowMeta, error) {
	return sw.storage.getWindowsMetadataByKey(cpId, key)
}

func (sw *storageWrapper[T]) closeWindow(cpId string, key string, winId string, watermark time.Duration, onFlush func(items []Message[T])) error {
	err := sw.storage.closeWindow(cpId, key, winId)
	if err != nil {
		return err
	}

	sw.taskInner.wg.Add(1)
	go func() {
		defer sw.taskInner.wg.Done()

		time.Sleep(watermark)

		items, err := sw.flushWindow(cpId, key, winId)
		if err != nil {
			panic(err)
		}

		onFlush(items)
	}()

	return nil
}

func (sw *storageWrapper[T]) flushWindow(cpId string, key string, winId string) ([]Message[T], error) {
	m, err := sw.storage.flushWindow(cpId, key, winId)
	if err != nil {
		return nil, err
	}

	msgs := make([]Message[T], len(m))

	for i := range m {
		msgs[i] = m[i].(Message[T])
	}

	return msgs, nil
}

func (sw *storageWrapper[T]) destroyWindow(cpId string, key string, winId string) error {
	return sw.storage.destroyWindow(cpId, key, winId)
}

func (sw *storageWrapper[T]) getWindowSize(cpId string, key string, winId string) (int, error) {
	return sw.storage.getWindowSize(cpId, key, winId)
}

func (sw *storageWrapper[T]) getKeys(cpId string) ([]string, error) {
	return sw.storage.getKeys(cpId)
}
