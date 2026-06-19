package internal

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

func GetTags(path string) ([]CTag, error) {
	cmd := exec.Command(
		"ctags",
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
