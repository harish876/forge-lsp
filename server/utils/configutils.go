package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func GetSectionNameFromUri(uri string) string {

	parts := strings.Split(uri, "/")
	if len(parts) == 0 {
		return ""
	}
	filename := parts[len(parts)-1]

	// Remove the ".py" extension
	filename = strings.TrimSuffix(filename, ".py")

	// Remove the "_job" suffix
	filename = strings.TrimSuffix(filename, "_job")

	return filename
}

// getLineContent retrieves the content of the specified line from the file
func GetLineContent(code []byte, lineNumber int) (string, error) {
	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(bytes.NewReader(code))
	currentLine := 0

	// Iterate through the lines of the file
	for scanner.Scan() {
		if currentLine == lineNumber {
			return scanner.Text(), nil
		}
		currentLine++
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("line %d does not exist in the file", lineNumber)
}
