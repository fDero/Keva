package history

import (
	"testing"

	"github.com/fDero/keva/misc"
)

func TestHistoryFileCreation(t *testing.T) {
	var errors [4]error
	var extracted_header historyFileHeader
	var hf HistoryFile
	pm := misc.NewMockPersistenceHandler()
	default_header := GetDefaultHistoryFileHeader(10, 20)
	hf, errors[0] = NewHistoryFile(pm, default_header)
	errors[1] = hf.AppendEvent("E1")
	errors[2] = hf.AppendEvent("E2")
	extracted_header, errors[3] = readHeader(pm)
	if err := misc.FirstOfManyErrorsOrNone(errors[:]); err != nil {
		t.Fatalf("Error during HistoryFile header extraction: %v", err)
	}
	if extracted_header.first_id != 10 {
		t.Fatalf("Expected first_id to be 10, got %d", extracted_header.first_id)
	}
	if extracted_header.entities_count != 2 {
		t.Fatalf("Expected entities_count to be 2, got %d", extracted_header.entities_count)
	}
	expected := misc.IterateValues("E1", "E2")
	gotten := hf.Iterate()
	for expected, gotten := range misc.Zip(expected, gotten) {
		if expected != gotten {
			t.Fatalf("Expected %s, got %s", expected, gotten)
		}
	}
}

func TestHistoryFileRestore(t *testing.T) {
	var errors [4]error
	var hf HistoryFile
	pm := misc.NewMockPersistenceHandler()
	pm.InitializeResource()
	errors[0] = writeHeader(pm, GetDefaultHistoryFileHeader(5, 100))
	errors[1] = pm.WriteBothLengthAndStringAtIndex(header_size_bytes, "E1")
	errors[2] = pm.WriteBothLengthAndStringAtIndex(header_size_bytes+8+2, "E2")
	hf, errors[3] = NewHistoryFile(pm, GetDefaultHistoryFileHeader(0, 0))
	if err := misc.FirstOfManyErrorsOrNone(errors[:]); err != nil {
		t.Fatalf("Error during HistoryFile restoration: %v", err)
	}
	if hf.header.first_id != 5 {
		t.Fatalf("Expected first_id to be 5, got %d", hf.header.first_id)
	}
	if hf.header.entities_count != 2 {
		t.Fatalf("Expected entities_count to be 2, got %d", hf.header.entities_count)
	}
}
