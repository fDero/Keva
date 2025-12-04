package core

type StorageSettings struct {
	internalStorage map[string]string
}

func NewStorageSettings() *StorageSettings {
	return &StorageSettings{
		internalStorage: make(map[string]string),
	}
}

func (ss *StorageSettings) ProcessEvent(event_raw_text string) error {
	return nil
}

func (ss *StorageSettings) DeleteRecord(key string) {
	delete(ss.internalStorage, key)
}

func (ss *StorageSettings) UpsertRecord(key string, value string) {
	ss.internalStorage[key] = value
}

func (ss *StorageSettings) FetchRecord(key string) (string, bool) {
	value, present := ss.internalStorage[key]
	return value, present
}
