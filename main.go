package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type RegexSlice []*regexp.Regexp

func (r *RegexSlice) String() string {
	patterns := make([]string, len(*r))
	for i, regex := range *r {
		patterns[i] = regex.String()
	}
	return strings.Join(patterns, ",")
}

func (r *RegexSlice) Set(value string) error {
	compiled, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regex pattern '%s': %v", value, err)
	}
	*r = append(*r, compiled)
	return nil
}

func main() {
	// Command line flags for the regex patterns - support multiple patterns
	var successPatterns RegexSlice
	var failurePatterns RegexSlice

	flag.Var(&successPatterns, "success", "Regex pattern that causes exit with status 0 (can be specified multiple times)")
	flag.Var(&failurePatterns, "failure", "Regex pattern that causes exit with status 1 (can be specified multiple times)")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := processLine(line)

		fmt.Println(processedLine)

		for _, pattern := range successPatterns {
			if pattern.MatchString(processedLine) {
				os.Exit(0)
			}
		}

		for _, pattern := range failurePatterns {
			if pattern.MatchString(processedLine) {
				os.Exit(1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(2)
	}
}

// processLine deals with control characters in the line
func processLine(line string) string {
	var result []rune

	for _, char := range line {
		switch {
		case char == '\b':
			// Remove the backspace and the previous character (if any)
			if len(result) > 0 {
				result = result[:len(result)-1]
			}

		case unicode.IsControl(char):
			// Skip all other control characters
			continue

		default:
			result = append(result, char)
		}
	}

	return string(result)
}
