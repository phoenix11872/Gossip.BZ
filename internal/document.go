package chunking

import (
	"fmt"
	"os"
)

type Document struct {
	DocumentUid string
	Content     string
}

// IngestFile reads the file at location without normalizing its contents, so
// whitespace (including spaces, newlines, and tabs) is retained exactly.
func IngestFile(location string) (Document, error) {
	content, err := os.ReadFile(location)
	if err != nil {
		return Document{}, fmt.Errorf("read document %q: %w", location, err)
	}

	return Document{
		DocumentUid: location,
		Content:     string(content),
	}, nil
}
