package misc

import (
	"errors"
	"os"
	"path/filepath"
)

func NewDiskPersistenceHandler(workdir, filename string) PersistenceHandler {
	fullpath := filepath.Join(workdir, filename)
	return PersistenceHandler{
		writeAtIndex: func(index int64, data []byte) error {
			return writeAtIndexToFile(fullpath, index, data)
		},
		readAtIndex: func(index int64, length int) ([]byte, error) {
			return readAtIndexFromFile(fullpath, index, length)
		},
		InitializeResource: func() error {
			return initializeResourceAsFile(fullpath)
		},
		ResourceExists: func() bool {
			return resourceExistsAsFile(fullpath)
		},
	}
}

func readAtIndexFromFile(fullpath string, index int64, length int) ([]byte, error) {
	f, err := os.OpenFile(fullpath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(index, 0)
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return []byte{}, errors.New("cannot read non-positive length")
	}
	buf := make([]byte, length)
	_, err = f.Read(buf)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

func resourceExistsAsFile(fullpath string) bool {
	return !os.IsNotExist(DiscardReturn(os.Stat(fullpath)))
}

func initializeResourceAsFile(fullpath string) error {
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func writeAtIndexToFile(fullpath string, index int64, data []byte) error {
	f, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return FirstOfManyErrorsOrNone([]error{
		DiscardReturn(f.Seek(index, 0)),
		DiscardReturn(f.Write(data)),
	})
}
