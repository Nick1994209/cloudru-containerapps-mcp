package utils

// deepMerge performs a deep merge of two maps
func DeepMerge(newData, oldData map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Copy old data
	for k, v := range oldData {
		merged[k] = v
	}

	// Merge new data (new data takes precedence, unless it's a zero value)
	for k, v := range newData {
		// If both values are maps, recursively merge them
		if oldVal, exists := merged[k]; exists {
			if oldMap, ok := oldVal.(map[string]interface{}); ok {
				if newMap, ok := v.(map[string]interface{}); ok {
					merged[k] = DeepMerge(newMap, oldMap)
					continue
				}
			}
			// If both values are arrays, merge each element individually
			if oldArray, ok := oldVal.([]interface{}); ok {
				if newArray, ok := v.([]interface{}); ok {
					// If arrays have the same length, merge each element
					if len(newArray) == len(oldArray) {
						mergedArray := make([]interface{}, len(newArray))
						for i := 0; i < len(newArray); i++ {
							// If both elements are maps, merge them
							if oldElem, ok := oldArray[i].(map[string]interface{}); ok {
								if newElem, ok := newArray[i].(map[string]interface{}); ok {
									mergedArray[i] = DeepMerge(newElem, oldElem)
								} else {
									mergedArray[i] = newArray[i]
								}
							} else {
								mergedArray[i] = newArray[i]
							}
						}
						merged[k] = mergedArray
						continue
					}
					// If arrays have different lengths, use the new array
					// This handles cases where the structure changes
					merged[k] = v
					continue
				}
			}
		}
		// If new value is zero, preserve old value
		if IsZeroValue(v) {
			continue
		}
		// Otherwise, new value takes precedence
		merged[k] = v
	}

	return merged
}

// IsZeroValue checks if a value is considered "zero" (empty string, 0, false)
// This is used to determine if a new value should be ignored in favor of the old value
// Note: nil is NOT considered a zero value here - it will be preserved
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case int64:
		return val == 0
	case float64:
		return val == 0
	case bool:
		return !val
	default:
		return false
	}
}
