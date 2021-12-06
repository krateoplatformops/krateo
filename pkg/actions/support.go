package actions

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path, prefix string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		if len(prefix) > 0 {
			ln = fmt.Sprintf("%s.%s", prefix, ln)
		}
		lines = append(lines, ln)
	}
	return lines, scanner.Err()
}

func printLines(filePath string, lines []string) error {
	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, ln := range lines {
		fmt.Fprintln(fp, ln)
	}

	return nil
}
