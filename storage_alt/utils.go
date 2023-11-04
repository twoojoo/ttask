package storage

func mergeMetadata(m1 map[string]int64, m2 map[string]int64) map[string]int64 {
	for k := range m2 {
		m1[k] = m2[k]
	}

	return m1
}
