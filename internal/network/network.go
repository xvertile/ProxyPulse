package network

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type NetworkStats struct {
	Interface string
	RXBytes   int
	TXBytes   int
}

func getNetworkStats() ([]NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var stats []NetworkStats

	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			return nil, scanner.Err()
		}
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}
		rxBytes, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		txBytes, err := strconv.Atoi(fields[9])
		if err != nil {
			return nil, err
		}
		stats = append(stats, NetworkStats{
			Interface: strings.TrimSuffix(fields[0], ":"),
			RXBytes:   rxBytes,
			TXBytes:   txBytes,
		})
	}

	return stats, scanner.Err()
}

func CalculateTransferRate(interval time.Duration) (int, int, error) {
	initialStats, err := getNetworkStats()
	if err != nil {
		return 0, 0, err
	}

	time.Sleep(interval)

	finalStats, err := getNetworkStats()
	if err != nil {
		return 0, 0, err
	}

	var totalRXBytes, totalTXBytes int

	for i := range initialStats {
		if initialStats[i].Interface == finalStats[i].Interface {
			totalRXBytes += finalStats[i].RXBytes - initialStats[i].RXBytes
			totalTXBytes += finalStats[i].TXBytes - initialStats[i].TXBytes
		}
	}

	totalRXMBps := math.Ceil((float64(totalRXBytes) * 8) / (float64(interval.Seconds()) * 1024 * 1024))
	totalTXMBps := math.Ceil((float64(totalTXBytes) * 8) / (float64(interval.Seconds()) * 1024 * 1024))

	return int(totalRXMBps), int(totalTXMBps), nil
}
