package utils

import (
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
