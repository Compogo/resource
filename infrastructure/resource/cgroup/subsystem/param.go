package subsystem

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/galsondor/go-ascii"
	"github.com/spf13/cast"
)

const (
	MaxTextValue = "max"
)

func ReadParamInt64(paramFilePath string) (int64, error) {
	param, err := ReadParam(paramFilePath)
	if err != nil {
		return 0, err
	}

	if param == "" || param == MaxTextValue {
		return 0, nil
	}

	return cast.ToInt64E(param)
}

func ReadParamUint64(paramFilePath string) (uint64, error) {
	param, err := ReadParam(paramFilePath)
	if err != nil {
		return 0, err
	}

	if param == "" || param == MaxTextValue {
		return 0, nil
	}

	return cast.ToUint64E(param)
}

func ReadParam(paramFilePath string) (string, error) {
	paramFile, err := os.OpenFile(paramFilePath, os.O_RDONLY, 0)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("cgroup.param.read: open param file '%s' failed: %w", paramFilePath, err)
	}

	if os.IsNotExist(err) {
		return "", nil
	}

	defer paramFile.Close()

	param, err := io.ReadAll(paramFile)
	if err != nil {
		return "", fmt.Errorf("cgroup.param.read: unable to read param file '%s': %w", paramFilePath, err)
	}

	return strings.TrimSpace(strings.Trim(string(param), string(rune(ascii.LF)))), nil
}
