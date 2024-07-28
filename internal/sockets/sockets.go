package sockets

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func GetTotalOpenSockets() (int32, error) {
	cmd := exec.Command("ss", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "TCP:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				totalTCPSockets, err := strconv.Atoi(fields[1])
				if err != nil {
					return 0, err
				}
				return int32(totalTCPSockets), nil
			}
		}
	}

	return 0, nil
}
