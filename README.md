# handytail

A Go utility for processing some input with pattern matching that results in success/failure exit codes, depending on what was found.

## Overview

`handytail` reads input line by line, processes control characters (discarding nearly all but backspace characters, which manipulate the line), and exits with specific status codes when certain regex patterns are matched. This makes it useful for watching command output or reading a file and determining success or failure conditions based on the contents rather than exit code or EOF.

## Features

- **Line-by-line processing** of stdin or file-based input, similar to grep
- **Backspace handling**: Removes backspace characters and the preceding character
- **Control character filtering**: Removes carriage returns and other control characters
- **Multiple regex patterns**: Support for multiple success and failure patterns
- **Flexible exit codes**:
  - 0 for success pattern matches
  - 1 for failure pattern matches
  - 2 for failures reading input
  - 3 for EOF

## Installation

```bash
go build -o handytail .
```

## Usage

```bash
./handytail [options] < input
command | ./handytail [options]
./handytail [options] /path/to/inputfile
```

### Command Line Options

- `-success <pattern>`: Regex pattern that causes exit with status 0 (can be specified multiple times)
- `-failure <pattern>`: Regex pattern that causes exit with status 1 (can be specified multiple times)
- `-quiet`: Disables printing the processed lines to screen
- `arg`: When provided, the program will read from this file instead of stdin


## Character Processing

### Backspace Handling
Input characters are processed as follows:
- `\b` (backspace): Removes the backspace and the previous character (if any)
- Multiple backspaces remove multiple preceding characters

```
Input:  "hello\b world"
Output: "hell world"

Input:  "test\b\b\b\bing"  
Output: "ing"
```

### Control Character Removal
- `\r` (carriage return): Completely removed
- `\t` (tab): Removed  
- `\n` (newline): Removed
- `\f` (form feed): Removed
- `\v` (vertical tab): Removed
- `\a` (bell): Removed
- Escape sequences: Removed

### Unicode Support
Unicode characters are preserved and work correctly with backspace operations.

## Exit Codes

- **0**: Success pattern matched
- **1**: Failure pattern matched  
- **2**: Error reading from stdin
- **No exit**: Program continues until EOF if no patterns match

## Real-World Examples

### CI/CD Pipeline Monitoring

```bash
# Monitor deployment script
./deploy.sh 2>&1 | ./handytail \
  -success "Deployment.*successful" \
  -success "All services.*running" \
  -failure "Deployment.*failed" \
  -failure "Service.*not.*responding" \
  -failure "Timeout.*exceeded"

if [ $? -eq 0 ]; then
    echo "Deployment succeeded!"
else
    echo "Deployment failed!"
fi
```

## Contributing

1. Write tests for new functionality
2. Ensure all existing tests pass
3. Update documentation as needed
4. Follow Go coding conventions
