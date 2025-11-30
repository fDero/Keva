package misc

import "encoding/binary"

type PersistenceHandler struct {
	writeAtIndex       func(index int64, data []byte) error
	readAtIndex        func(index int64, length int) ([]byte, error)
	ResourceExists     func() bool
	InitializeResource func() error
}

func (d *PersistenceHandler) ReadInt64AtIndex(index int64) (int64, int, error) {
	buf, err := d.readAtIndex(index, 8)
	if err != nil {
		return -1, 0, err
	}
	return int64(binary.LittleEndian.Uint64(buf)), 8, nil
}

func (d *PersistenceHandler) ReadInt32AtIndex(index int64) (int32, int, error) {
	buf, err := d.readAtIndex(index, 4)
	if err != nil {
		return -1, 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), 4, nil
}

func (d *PersistenceHandler) ReadStringAtIndex(index int64, length int) (string, int, error) {
	buffer, err := d.readAtIndex(index, length)
	if err != nil {
		return "", 0, err
	}
	return string(buffer), length, nil
}

func (d *PersistenceHandler) ReadBothLengthAndStringAtIndex(index int64) (string, int, error) {
	var errors [2]error
	var length int64 = 0
	var text string = ""
	length, _, errors[0] = d.ReadInt64AtIndex(index)
	text, _, errors[1] = d.ReadStringAtIndex(index+8, int(length))
	return text, 8 + int(length), FirstOfManyErrorsOrNone(errors[:])
}

func (d *PersistenceHandler) WriteInt64AtIndex(index int64, data int64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(data))
	return d.writeAtIndex(index, buf)
}

func (d *PersistenceHandler) WriteInt32AtIndex(index int64, data int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(data))
	return d.writeAtIndex(index, buf)
}

func (d *PersistenceHandler) WriteBothLengthAndStringAtIndex(index int64, data string) error {
	return FirstOfManyErrorsOrNone([]error{
		d.WriteInt64AtIndex(index, int64(len(data))),
		d.WriteStringAtIndex(index+8, data),
	})
}

func (d *PersistenceHandler) WriteStringAtIndex(index int64, data string) error {
	strBytes := []byte(data)
	buf := make([]byte, 4+len(strBytes))
	binary.LittleEndian.PutUint32(buf[:4], uint32(len(strBytes)))
	copy(buf[4:], strBytes)
	return d.writeAtIndex(index, buf)
}
