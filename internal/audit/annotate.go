package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Annotation holds a user-defined note attached to a secret path.
type Annotation struct {
	Path      string    `json:"path"`
	Note      string    `json:"note"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// AnnotationStore is a map of path -> Annotation.
type AnnotationStore map[string]Annotation

// SaveAnnotations writes the annotation store to a JSON file.
func SaveAnnotations(path string, store AnnotationStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal annotations: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadAnnotations reads an annotation store from a JSON file.
func LoadAnnotations(path string) (AnnotationStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read annotations: %w", err)
	}
	var store AnnotationStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal annotations: %w", err)
	}
	return store, nil
}

// ApplyAnnotations attaches notes to CompareResults where a matching annotation exists.
func ApplyAnnotations(reports []CompareResult, store AnnotationStore) []CompareResult {
	annotated := make([]CompareResult, len(reports))
	for i, r := range reports {
		if ann, ok := store[r.Path]; ok {
			r.Note = ann.Note
		}
		annotated[i] = r
	}
	return annotated
}
