package ttask

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
