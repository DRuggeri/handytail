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
	quiet := false
	showHelp := false

	flag.Var(&successPatterns, "success", "Regex pattern that causes exit with status 0 (can be specified multiple times)")
	flag.Var(&failurePatterns, "failure", "Regex pattern that causes exit with status 1 (can be specified multiple times)")
	flag.BoolVar(&quiet, "quiet", false, "Suppress output to stdout of the file contents")
	flag.BoolVar(&showHelp, "help", false, "Show help message and exit")

	flag.Usage = printHelp

	flag.Parse()

	if showHelp {
		printHelp()
	}

	args := flag.Args()

	input := os.Stdin
	if len(args) > 0 {
		if len(args) > 1 {
			fmt.Fprintf(os.Stderr, "Too many positional arguments provided.\n\n")
			printHelp()
		}
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file '%s': %v\n", args[0], err)
			os.Exit(2)
		}
		defer file.Close()
		input = file
	}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := processLine(line)

		if !quiet {
			fmt.Println(processedLine)
		}

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

	fmt.Fprintf(os.Stderr, "EOF found but no matching lines were detected")
	os.Exit(3)
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

// printHelp displays usage and exit code information
func printHelp() {
	fmt.Fprintf(os.Stderr, "handytail - process input with pattern matching and exit codes\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  ./handytail [options] < input\n")
	fmt.Fprintf(os.Stderr, "  command | ./handytail [options]\n")
	fmt.Fprintf(os.Stderr, "  ./handytail [options] /path/to/inputfile\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -success <pattern>   Regex pattern that causes exit with status 0 (can be specified multiple times)\n")
	fmt.Fprintf(os.Stderr, "  -failure <pattern>   Regex pattern that causes exit with status 1 (can be specified multiple times)\n")
	fmt.Fprintf(os.Stderr, "  -quiet              Disables printing the processed lines to screen\n")
	fmt.Fprintf(os.Stderr, "  -help               Show this help message and exit\n")
	fmt.Fprintf(os.Stderr, "  arg                 When provided, the program will read from this file instead of stdin\n\n")
	fmt.Fprintf(os.Stderr, "Exit Codes:\n")
	fmt.Fprintf(os.Stderr, "  0  Success pattern matched\n")
	fmt.Fprintf(os.Stderr, "  1  Failure pattern matched\n")
	fmt.Fprintf(os.Stderr, "  2  Error reading from stdin or opening file\n")
	fmt.Fprintf(os.Stderr, "  3  EOF found but no matching lines were detected\n")
	os.Exit(99)
}
