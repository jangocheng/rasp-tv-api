package data

import (
	"fmt"
	"os"
)

// IsProcessRunning checks if pid is running
func IsProcessRunning(pid int64) bool {
	pidStr := fmt.Sprintf("%d", pid)

	stat, err := os.Open("/proc")
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
