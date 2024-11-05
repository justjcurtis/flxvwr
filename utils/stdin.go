package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadStdin() ([]string, error) {
	file := os.Stdin
	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	size := fi.Size()
	if size > 0 {
		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) > 0 {
				lines = append(lines, line)
			}
		}
		return lines, nil
	}
	return nil, fmt.Errorf("no input provided")
}
