package core

import "errors"

type StorageSettings struct {
	internalStorage map[string]string
}

func NewStorageSettings() *StorageSettings {
	return &StorageSettings{
		internalStorage: make(map[string]string),
	}
}

func (ss *StorageSettings) ProcessEvent(event Event) error {
	callbacks := map[string]func(Event){
		"UPSERT/VALUE": ss.UpsertRecord,
		"DELETE/VALUE": ss.DeleteRecord,
	}
	callback, present := callbacks[event.action]
	if !present {
		return errors.New("unknown event action: [" + event.action + "] for key: [" + event.key + "]")
	}
	callback(event)
	return nil
}

func (ss *StorageSettings) DeleteRecord(event Event) {
	delete(ss.internalStorage, event.key)
}

func (ss *StorageSettings) UpsertRecord(event Event) {
	ss.internalStorage[event.key] = event.new_value
}

func (ss *StorageSettings) FetchRecord(key string) (string, bool) {
	value, present := ss.internalStorage[key]
	return value, present
}
