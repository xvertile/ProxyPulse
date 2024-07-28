package cpu

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetTotalCPUUsage(processName string) (float64, error) {
	pidCmd := exec.Command("pgrep", processName)
	pidOut, err := pidCmd.Output()
	if err != nil {
		return 0, err
	}

	pid := strings.TrimSpace(string(pidOut))
	if pid == "" {
		return 0, fmt.Errorf("process %s not found", processName)
	}

	initialTotalCPU, err := readTotalCPU()
	if err != nil {
		return 0, err
	}
	initialProcessCPU, err := readProcessCPU(pid)
	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)

	finalTotalCPU, err := readTotalCPU()
	if err != nil {
		return 0, err
	}
	finalProcessCPU, err := readProcessCPU(pid)
	if err != nil {
		return 0, err
	}

	totalCPUDelta := finalTotalCPU - initialTotalCPU
	processCPUDelta := finalProcessCPU - initialProcessCPU

	if totalCPUDelta == 0 {
		return 0, fmt.Errorf("total CPU time did not change")
	}

	cpuUsage := (float64(processCPUDelta) / float64(totalCPUDelta)) * 100
	return cpuUsage, nil
}

func readTotalCPU() (int64, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				return 0, fmt.Errorf("unexpected /proc/stat format")
			}
			var totalCPU int64
			for _, val := range fields[1:8] {
				cpuTime, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return 0, err
				}
				totalCPU += cpuTime
			}
			return totalCPU, nil
		}
	}
	return 0, fmt.Errorf("could not find total CPU usage in /proc/stat")
}

func readProcessCPU(pid string) (int64, error) {
	data, err := ioutil.ReadFile("/proc/" + pid + "/stat")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 17 {
		return 0, fmt.Errorf("unexpected /proc/[pid]/stat format")
	}

	utime, err := strconv.ParseInt(fields[13], 10, 64)
	if err != nil {
		return 0, err
	}
	stime, err := strconv.ParseInt(fields[14], 10, 64)
	if err != nil {
		return 0, err
	}

	return utime + stime, nil
}
