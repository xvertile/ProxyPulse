package memory

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetProcessMemoryUsage(processName string) (float32, error) {
	pidCmd := exec.Command("pgrep", processName)
	pidOut, err := pidCmd.Output()
	if err != nil {
		return 0, err
	}

	pid := strings.TrimSpace(string(pidOut))
	if pid == "" {
		return 0, fmt.Errorf("process %s not found", processName)
	}

	file, err := os.Open("/proc/" + pid + "/status")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			memoryUsageKB, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			memoryUsageGB := float32(memoryUsageKB) / (1024 * 1024)
			return memoryUsageGB, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("memory usage information not found for process %s", processName)
}
