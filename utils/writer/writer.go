package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

func ToFile(filePath string, bytes []byte) (int, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return 0, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writtenBytesNumber, err := file.Write(bytes)
	if err != nil {
		return 0, err
	}

	if len(bytes) > 0 && writtenBytesNumber == 0 {
		return 0, fmt.Errorf("unable to write bytes into file: nothing was written")
	}

	return writtenBytesNumber, nil
}
