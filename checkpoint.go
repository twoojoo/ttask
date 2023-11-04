package ttask

// Turns an opertator into a checkpoint.
//
// If the processing of a message is interrupted before reaching the next checkpoint (windows act as checkpoints) or before the last
// task operation is exectuted, the next task execution will recover the checkpointed messages.
//
// The recovery procedure will start when a new message reach the checkpoint.
//
// DON'T USE if you have more than one replica of the same task and you're using a remote storage like Redis,
// unless you find a way to set a different task id for each replica and you're able to retain it on each restart.
//
// USELESS if your using the in-memory storage or if you're using a local storage you're running the task
// in a volatile resource (e.g. K8s pods)
func Checkpoint[T, R any](id string, operator Operator[T, R]) Operator[T, R] {
	first := true

	return func(inner *Inner, x *Message[T], next *Step) {
		if first {
			err := recoverCheckpoint(inner.storage, inner.TaskID(), id, func(m *Message[T]) {
				operator(inner, m, next)
				inner.storage.clearCheckpoint(inner.TaskID(), m.Id, id)
			})

			if err != nil {
				inner.Error(err)
			}

			first = false
		}

		storeCheckpoint(inner.storage, inner.TaskID(), id, x)

		operator(inner, x, next)

		inner.storage.clearCheckpoint(inner.TaskID(), x.Id, id)
	}
}
