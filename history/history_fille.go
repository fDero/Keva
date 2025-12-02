package history

import (
	"io"
	"iter"

	"github.com/fDero/keva/misc"
)

func GetDefaultHistoryFileHeader(first_id, current_epoch int64) historyFileHeader {
	return historyFileHeader{
		first_id:       first_id,
		entities_count: 0,
		size_mem_bytes: header_size_bytes,
		creation_epoch: current_epoch,
		protocol_major: 1,
		protocol_minor: 0,
		protocol_patch: 0,
	}
}

func NewHistoryFile(pm misc.PersistenceHandler, default_header historyFileHeader) (HistoryFile, error) {
	var errors [2]error
	var file HistoryFile
	file.manager = pm
	if !pm.ResourceExists() {
		pm.InitializeResource()
		writeHeader(pm, default_header)
	}
	file.header, errors[0] = readHeader(pm)
	file.content, errors[1] = restoreHistoryFileContents(pm, file.header.first_id)
	file.header.entities_count = int64(len(file.content))
	return file, misc.FirstOfManyErrorsOrNone(errors[:])
}

func restoreHistoryFileContents(pm misc.PersistenceHandler, first_id int64) (map[int64]string, error) {
	var cursor int64 = header_size_bytes
	var counter int64 = 0
	content := make(map[int64]string)
	for {
		event_text, offset, err := pm.ReadBothLengthAndStringAtIndex(cursor)
		if err == io.EOF {
			return content, nil
		}
		if err != nil {
			return nil, err
		}
		content[first_id+counter] = event_text
		cursor += int64(offset)
		counter++
	}
}

func (hf *HistoryFile) GetEventByID(id int64) string {
	return hf.content[id]
}

func (hf *HistoryFile) GetLastEventID() int64 {
	return hf.header.first_id + hf.header.entities_count - 1
}

func (hf *HistoryFile) AppendEvent(event_text string) error {
	len64 := int64(len(event_text))
	content_index := hf.header.first_id + hf.header.entities_count
	hf.content[content_index] = event_text
	err := hf.manager.WriteBothLengthAndStringAtIndex(hf.header.size_mem_bytes, event_text)
	if err != nil {
		delete(hf.content, content_index)
		return err
	}
	hf.header.entities_count++
	hf.header.size_mem_bytes += int64(8 + len64)
	return writeHeader(hf.manager, hf.header)
}

func (hf *HistoryFile) Iterate() iter.Seq[string] {
	first_id := hf.header.first_id
	last_id := first_id + hf.header.entities_count
	return misc.IterateMapValues(hf.content, misc.Range(first_id, last_id+1))
}
