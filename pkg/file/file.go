package file

import (
	"fmt"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("file: no path provided")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file: cannot open file - %v", err.Error())
	}

	return data, nil
}
