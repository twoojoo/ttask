package ttask

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error {
	return nil
}

func (s *MemoryStorage) clearCheckpoint(taskId string, msgId string, cpId string) error {
	return nil
}

func (s *MemoryStorage) getCheckpointMessages(taskId string, cpId string) ([][]byte, error) {
	return [][]byte{}, nil
}
