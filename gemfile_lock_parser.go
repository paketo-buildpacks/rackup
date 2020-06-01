package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GemfileLockParser struct{}

func NewGemfileLockParser() GemfileLockParser {
	return GemfileLockParser{}
}

func (p GemfileLockParser) Parse(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}

		return false, fmt.Errorf("failed to parse Gemfile.lock: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "GEM" {
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "rack") {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
