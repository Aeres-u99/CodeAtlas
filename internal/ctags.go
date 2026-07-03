package internal

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func GetCtagsPath() string {
	if v, ok := os.LookupEnv("HERMES_CTAGS"); ok {
		slog.Info("Found Custom CTAGS: %s\n", v)
		return v
	} else {
		return "ctags"
	}
}

var cmdBin = GetCtagsPath()

func GetTags(path string) ([]CTag, error) {
	cmd := exec.Command(
		cmdBin,
		"--output-format=json",
		"--fields=+na",
		path,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(stdout.String(), "\n")

	var tags []CTag

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var tag CTag

		if err := json.Unmarshal([]byte(line), &tag); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
