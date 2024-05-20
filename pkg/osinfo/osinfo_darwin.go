package osinfo

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

var defaultReturnValue = OSInfo{
	OS:       "darwin",
	Version:  "unknown",
	Platform: "unknown",
}

func GetOSInfo() OSInfo {
	out, err := getInfo()
	if err != nil {
		return defaultReturnValue
	}
	for strings.Contains(out, brokenPipeOutput) {
		out, err = getInfo()
		time.Sleep(retryWaitTime)
	}
	if err != nil {
		return defaultReturnValue
	}
	osStr := replaceLineTerminations(out)
	osInfo := strings.Split(osStr, " ")
	if len(osInfo) != expectedOSInfoArrLen {
		return defaultReturnValue
	}
	oss := OSInfo{
		OS:       osInfo[0],
		Version:  osInfo[1],
		Platform: osInfo[2],
	}
	return oss
}

func getInfo() (string, error) {
	cmd := exec.Command("uname", "-srm")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), err
}
