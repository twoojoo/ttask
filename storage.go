package ttask

type Storage interface {
	storeCheckpoint(taskId string, msgId string, cpId string, msg []byte) error
	clearCheckpoint(taskId string, msgId string, cpId string) error
	getCheckpointMessages(taskId string, cpId string) ([][]byte, error)
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
