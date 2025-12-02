package history

import "github.com/fDero/keva/misc"

const (
	size_mem_bytes_offset = 0
	entities_count_offset = 8
	creation_epoch_offset = 16
	protocol_major_offset = 24
	protocol_minor_offset = 28
	protocol_patch_offset = 32
	first_id_offset       = 36
	header_size_bytes     = 44
)

func writeHeader(pm misc.PersistenceHandler, h historyFileHeader) error {
	var errors [7]error
	errors[0] = pm.WriteInt64AtIndex(size_mem_bytes_offset, h.size_mem_bytes)
	errors[1] = pm.WriteInt64AtIndex(entities_count_offset, h.entities_count)
	errors[2] = pm.WriteInt64AtIndex(creation_epoch_offset, h.creation_epoch)
	errors[3] = pm.WriteInt32AtIndex(protocol_major_offset, h.protocol_major)
	errors[4] = pm.WriteInt32AtIndex(protocol_minor_offset, h.protocol_minor)
	errors[5] = pm.WriteInt32AtIndex(protocol_patch_offset, h.protocol_patch)
	errors[6] = pm.WriteInt64AtIndex(first_id_offset, h.first_id)
	return misc.FirstOfManyErrorsOrNone(errors[:])
}

func readHeader(pm misc.PersistenceHandler) (historyFileHeader, error) {
	var errors [7]error
	var header historyFileHeader
	header.size_mem_bytes, _, errors[0] = pm.ReadInt64AtIndex(size_mem_bytes_offset)
	header.entities_count, _, errors[1] = pm.ReadInt64AtIndex(entities_count_offset)
	header.creation_epoch, _, errors[2] = pm.ReadInt64AtIndex(creation_epoch_offset)
	header.protocol_major, _, errors[3] = pm.ReadInt32AtIndex(protocol_major_offset)
	header.protocol_minor, _, errors[4] = pm.ReadInt32AtIndex(protocol_minor_offset)
	header.protocol_patch, _, errors[5] = pm.ReadInt32AtIndex(protocol_patch_offset)
	header.first_id, _, errors[6] = pm.ReadInt64AtIndex(first_id_offset)
	return header, misc.FirstOfManyErrorsOrNone(errors[:])
}
