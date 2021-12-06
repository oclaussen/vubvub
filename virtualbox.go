package vubvub

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"time"
)

func vbm(args ...string) (string, error) {
	return vbmRetry(5, args...)
}

func vbmRetry(retry int, args ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	vboxManage, err := exec.LookPath("VBoxManage")
	if err != nil {
		return "", errors.New("could not find VBoxManage command")
	}

	cmd := exec.Command(vboxManage, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}

	if retry > 1 {
		if strings.Contains(stderr.String(), "error: The object is not ready") {
			time.Sleep(100 * time.Millisecond)

			return vbmRetry(retry-1, args...)
		}
	}

	return stdout.String(), nil
}
