package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type OwnerEntry struct {
	Path    string `json:"path"`
	Owner   string `json:"owner"`
	Team    string `json:"team"`
	Contact string `json:"contact"`
}

type OwnershipStore struct {
	Owners []OwnerEntry `json:"owners"`
}

func SaveOwnership(path string, store OwnershipStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal ownership: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadOwnership(path string) (OwnershipStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return OwnershipStore{}, nil
		}
		return OwnershipStore{}, fmt.Errorf("read ownership: %w", err)
	}
	var store OwnershipStore
	if err := json.Unmarshal(data, &store); err != nil {
		return OwnershipStore{}, fmt.Errorf("parse ownership: %w", err)
	}
	return store, nil
}

func ApplyOwnership(reports []CompareReport, store OwnershipStore) []CompareReport {
	for i, r := range reports {
		for _, o := range store.Owners {
			if strings.HasPrefix(r.Path, o.Path) {
				if reports[i].Annotations == nil {
					reports[i].Annotations = map[string]string{}
				}
				reports[i].Annotations["owner"] = o.Owner
				reports[i].Annotations["team"] = o.Team
				reports[i].Annotations["contact"] = o.Contact
				break
			}
		}
	}
	return reports
}

func LookupOwner(store OwnershipStore, path string) (OwnerEntry, bool) {
	for _, o := range store.Owners {
		if strings.HasPrefix(path, o.Path) {
			return o, true
		}
	}
	return OwnerEntry{}, false
}
