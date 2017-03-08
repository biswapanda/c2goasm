package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func isLower(str string) bool {

	for _, r := range str {
		return unicode.IsLower(r)
	}
	return false
}

func removeUndefined(line, undef string) string {

	if parts := strings.SplitN(line, undef, 2); len(parts) > 1 {
		line = parts[0] + strings.TrimSpace(parts[1])
	}
	return line
}

func process(lines []string) ([]string, error) {

	var result []string

	for _, line := range lines {

		// Remove ## comments
		if parts := strings.SplitN(line, `##`, 2); len(parts) > 1 {
			if strings.TrimSpace(parts[0])== "" {
				continue
			}
			line = parts[0]
		}

		// Skip lines with aligns
		if strings.Contains(line, ".align") {
			continue
		}

		// Make jmps uppercase
		if parts := strings.SplitN(line, `LBB0`, 2); len(parts) > 1 {
			// unless it is a label
			if !strings.Contains(parts[1], ":") {
				// make jmp statement uppercase
				line = strings.ToUpper(parts[0]) + "LBB0" + parts[1]
			}
		}

		fields := strings.Fields(line)
		// Test for any non-jmp instruction (lower case mnemonic)
		if len(fields) > 0 && !strings.Contains(fields[0], ":") && isLower(fields[0]) {
			// prepend line with comment for subsequent asm2plan9s assembly
			line = "                                 // " + strings.TrimSpace(line)
		}

		line = removeUndefined(line, "ptr")
		line = removeUndefined(line, "xmmword")
		line = removeUndefined(line, "ymmword")

		result = append(result, line)
	}

	return result, nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("error: no input specified\n\n")
		fmt.Println("usage: cagac file")
		fmt.Println("  will ")
		return
	}
	fmt.Println("Processing", os.Args[1])
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	result, err := process(lines)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	err = writeLines(result, os.Args[1])
	if err != nil {
		log.Fatalf("writeLines: %s", err)
	}
}