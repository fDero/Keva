package misc

import "io"

func NewMockPersistenceHandler() PersistenceHandler {
	var buffer []byte
	initialized := false

	return PersistenceHandler{
		writeAtIndex: func(index int64, data []byte) error {
			return writeAtIndexToBuffer(&buffer, index, data)
		},
		readAtIndex: func(index int64, length int) ([]byte, error) {
			return readAtIndexFromBuffer(&buffer, index, length)
		},
		InitializeResource: func() error {
			initialized = true
			buffer = make([]byte, 0)
			return nil
		},
		ResourceExists: func() bool {
			return initialized
		},
	}
}

func writeAtIndexToBuffer(buffer *[]byte, index int64, data []byte) error {
	requiredSize := int(index) + len(data)
	if len(*buffer) < requiredSize {
		newBuffer := make([]byte, requiredSize)
		copy(newBuffer, *buffer)
		*buffer = newBuffer
	}
	copy((*buffer)[index:], data)
	return nil
}

func readAtIndexFromBuffer(buffer *[]byte, index int64, length int) ([]byte, error) {
	if int(index) >= len(*buffer) {
		return nil, io.EOF
	}
	endIndex := int(index) + length
	if endIndex > len(*buffer) {
		endIndex = len(*buffer)
	}
	result := make([]byte, endIndex-int(index))
	copy(result, (*buffer)[index:endIndex])
	if len(result) < length {
		return result, io.EOF
	}
	return result, nil
}
