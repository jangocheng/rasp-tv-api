package data

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsProcessRunning(pid int64) bool {
	pidStr := fmt.Sprintf("%d", pid)
	statPath := filepath.Join("/proc", pidStr, "stat")

	stat, err := os.Open(statPath)
	if err != nil {
		return false
	}

	names, err := stat.Readdirnames(-1)
	if err != nil {
		return false
	}

	for _, name := range names {
		if name == pidStr {
			return true
		}
	}

	return false
}
