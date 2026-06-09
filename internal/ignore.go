package internal

import (
	"bufio"
	"fmt"
	ignore "github.com/denormal/go-gitignore"
	"os"
	"strings"
)

var DefaultIgnore = map[string]bool{
	".git":         true,
	".hermes":      true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	"target":       true,
}

func LoadIgnoreFile(root string) (map[string]bool, error) {
	ignore := make(map[string]bool)

	for k, v := range DefaultIgnore {
		ignore[k] = v
	}

	fmt.Println("Loading:", root+"/.hermesignore")
	file, err := os.Open(root + "/.hermesignore")
	if err != nil {
		if os.IsNotExist(err) {
			return ignore, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimSuffix(line, "/")
		fmt.Printf("Read line: [%s]\n", line)

		if line == "" {
			continue
		}
		ignore[line] = true
	}
	fmt.Println("Loaded ignore:", ignore)
	return ignore, scanner.Err()
}

func ShouldIgnore(name string, ignore map[string]bool) bool {
	name = strings.TrimSuffix(name, "/")
	return ignore[name]
}
