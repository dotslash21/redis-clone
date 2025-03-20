package config

import (
	"regexp"
	"strings"
	"sync"
)

type store struct {
	settings map[string]string
	mu       sync.RWMutex
}

var storeInstance *store = &store{
	settings: make(map[string]string),
	mu:       sync.RWMutex{},
}

// SetConfig sets a configuration value
func SetConfig(key, value string) {
	storeInstance.mu.Lock()
	defer storeInstance.mu.Unlock()
	storeInstance.settings[key] = value
}

// GetConfig retrieves a configuration value
func GetConfig(searchPattern string) (map[string]string, error) {
	storeInstance.mu.RLock()
	defer storeInstance.mu.RUnlock()

	results := make(map[string]string)

	// Convert Redis glob pattern to regex pattern
	regexPattern := convertGlobToRegex(searchPattern)

	// Compile the regex pattern
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	// Find all the matching keys
	for key, value := range storeInstance.settings {
		if regex.MatchString(key) {
			results[key] = value
		}
	}

	return results, nil
}

// Helper to convert Redis glob pattern to regex
func convertGlobToRegex(pattern string) string {
	// First escape the entire string to handle special regex chars
	escaped := regexp.QuoteMeta(pattern)

	// Then replace escaped wildcards with their regex equivalents
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = strings.ReplaceAll(escaped, "\\?", ".")

	// Ensure we match the entire string
	return "^" + escaped + "$"
}
