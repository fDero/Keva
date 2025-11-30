package misc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessMetadata struct {
	Pid         int    `json:"pid"`
	StartTime   string `json:"start_time"`
	ProcessName string `json:"process_name"`
}

func NewProcessMetadata(pid int) (*ProcessMetadata, error) {
	pid32 := int32(pid)
	p, err := process.NewProcess(pid32)
	if err != nil {
		return nil, err
	}
	start_time, _ := p.CreateTime()
	name, _ := p.Name()
	return &ProcessMetadata{
		Pid:         pid,
		StartTime:   fmt.Sprintf("%d", start_time),
		ProcessName: name,
	}, nil
}

func CreateLockFile(file_path string) error {
	metadata, err := NewProcessMetadata(os.Getpid())
	if err != nil {
		fmt.Println("Error creating process metadata:", err)
		return err
	}
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := os.WriteFile(file_path, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	fmt.Println("Saved metadata to: ", file_path)
	return nil
}

func PurgeUnusedLockFiles(file_path string) error {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return nil
	}
	var old_process_metadata ProcessMetadata
	if err := json.Unmarshal(data, &old_process_metadata); err != nil {
		fmt.Println("Error decoding lock file:", err)
		return err
	}
	_, err = NewProcessMetadata(old_process_metadata.Pid)
	if err != nil {
		fmt.Println("Removing stale lock file:", file_path)
		os.Remove(file_path)
		return nil
	}
	pid_str := fmt.Sprintf("%d", old_process_metadata.Pid)
	panic("Another instance is already running in the same workdir with PID: " + pid_str)
}
