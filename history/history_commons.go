package history

import "github.com/fDero/keva/misc"

type historyFileHeader struct {
	size_mem_bytes int64
	entities_count int64
	creation_epoch int64
	protocol_major int32
	protocol_minor int32
	protocol_patch int32
	first_id       int64
}

type HistoryFile struct {
	manager misc.PersistenceHandler
	header  historyFileHeader
	content map[int64]string
}
