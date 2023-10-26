package utils

import (
	"os"
)

func OpenOrCreateFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return file, nil
}
