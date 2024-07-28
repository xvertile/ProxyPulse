package filedescriptors

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func GetTotalFileDescriptors(processName string) (int32, error) {
	pidCmd := exec.Command("pgrep", processName)
	var pidOut bytes.Buffer
	pidCmd.Stdout = &pidOut

	err := pidCmd.Run()
	if err != nil {
		return 0, err
	}

	pid := strings.TrimSpace(pidOut.String())
	if pid == "" {
		return 0, nil
	}

	fdCmd := exec.Command("ls", "/proc/"+pid+"/fd")
	var fdOut bytes.Buffer
	fdCmd.Stdout = &fdOut

	err = fdCmd.Run()
	if err != nil {
		return 0, err
	}

	fdCountCmd := exec.Command("wc", "-l")
	fdCountCmd.Stdin = &fdOut
	var countOut bytes.Buffer
	fdCountCmd.Stdout = &countOut

	err = fdCountCmd.Run()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(countOut.String()))
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}
