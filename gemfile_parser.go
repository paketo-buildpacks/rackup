package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type GemfileParser struct{}

func NewGemfileParser() GemfileParser {
	return GemfileParser{}
}

func (p GemfileParser) Parse(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to parse Gemfile: %w", err)
	}
	defer file.Close()

	mriRe := regexp.MustCompile(`^ruby .*`)
	hasMri := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		if hasMri == false {
			hasMri = mriRe.Match(line)
			if hasMri == true {
				return true, nil
			}
		}
	}

	return hasMri, nil
}
