package utils

import (
	"testing"
)

func TestDeepMerge(t *testing.T) {
	tests := []struct {
		name     string
		newData  map[string]interface{}
		oldData  map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "simple merge - new data overrides old",
			newData: map[string]interface{}{
				"key1": "value1_new",
				"key2": "value2_new",
				"key3": map[string]interface{}{
					"key4": "",
					"key5": "value2_new",
					"key6": 10,
					"key7": 0,
				},
			},
			oldData: map[string]interface{}{
				"key1": "value1_old",
				"key2": "value2_old",
				"key3": map[string]interface{}{
					"key4": "value4_old",
					"key5": "value5_old",
					"key6": 6,
					"key7": 7,
				},
				"key4": "value3_old",
			},
			expected: map[string]interface{}{
				"key1": "value1_new",
				"key2": "value2_new",
				"key3": map[string]interface{}{
					"key4": "value4_old",
					"key5": "value2_new",
					"key6": 10,
					"key7": 7,
				},
				"key4": "value3_old",
			},
		},
		{
			name: "nested map merge",
			newData: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"key1": "value1_new",
						"key2": "value2_new",
					},
				},
			},
			oldData: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"key1": "value1_old",
						"key3": "value3_old",
					},
				},
				"key4": "value4_old",
			},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"key1": "value1_new",
						"key2": "value2_new",
						"key3": "value3_old",
					},
				},
				"key4": "value4_old",
			},
		},
		{
			name:    "empty new data",
			newData: map[string]interface{}{},
			oldData: map[string]interface{}{
				"key1": "value1_old",
				"key2": "value2_old",
			},
			expected: map[string]interface{}{
				"key1": "value1_old",
				"key2": "value2_old",
			},
		},
		{
			name: "empty old data",
			newData: map[string]interface{}{
				"key1": "value1_new",
				"key2": "value2_new",
			},
			oldData: map[string]interface{}{},
			expected: map[string]interface{}{
				"key1": "value1_new",
				"key2": "value2_new",
			},
		},
		{
			name: "deeply nested merge",
			newData: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "value_d_new",
						},
					},
				},
			},
			oldData: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "value_d_old",
							"e": "value_e_old",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "value_d_new",
							"e": "value_e_old",
						},
					},
				},
			},
		},
		{
			name: "mixed types at same level",
			newData: map[string]interface{}{
				"string": "new_string",
				"number": 42,
				"bool":   true,
			},
			oldData: map[string]interface{}{
				"string": "old_string",
				"number": 100,
				"bool":   false,
			},
			expected: map[string]interface{}{
				"string": "new_string",
				"number": 42,
				"bool":   true,
			},
		},
		{
			name: "arrays are not merged (new replaces old)",
			newData: map[string]interface{}{
				"array": []interface{}{"new1", "new2"},
			},
			oldData: map[string]interface{}{
				"array": []interface{}{"old1", "old2"},
			},
			expected: map[string]interface{}{
				"array": []interface{}{"new1", "new2"},
			},
		},
		{
			name: "nil values",
			newData: map[string]interface{}{
				"key1": nil,
				"key2": "value2",
			},
			oldData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2_old",
			},
			expected: map[string]interface{}{
				"key1": nil,
				"key2": "value2",
			},
		},
		{
			name: "array with nested objects merge",
			newData: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{
						"key2": "new_value2",
						"key4": map[string]interface{}{
							"key5": "new_key_5",
						},
					},
				},
			},
			oldData: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{
						"key2": "new_value2",
						"key3": "old_value3",
						"key4": map[string]interface{}{
							"key5": "old_key_5",
							"key6": "old_key_6",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{
						"key2": "new_value2",
						"key3": "old_value3",
						"key4": map[string]interface{}{
							"key5": "new_key_5",
							"key6": "old_key_6",
						},
					},
				},
			},
		},
		{
			name: "job template with resources merge",
			newData: map[string]interface{}{
				"template": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"resources": map[string]interface{}{
								"cpu":    "0.3",
								"memory": "768Mi",
							},
						},
					},
				},
			},
			oldData: map[string]interface{}{
				"name": "test-job-min-130045",
				"template": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"image": "nvkorolkov-public.cr.cloud.ru/job-did-you-know@sha256:292dd73a2cf83cc0d27bbac1e14a00c71fcb816cbb4f63d3f5e95be685171fa8",
							"name":  "test-job-min-130045",
							"resources": map[string]interface{}{
								"cpu":    "0.1",
								"memory": "256Mi",
							},
							"volumeMounts": []interface{}{},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"name": "test-job-min-130045",
				"template": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"image": "nvkorolkov-public.cr.cloud.ru/job-did-you-know@sha256:292dd73a2cf83cc0d27bbac1e14a00c71fcb816cbb4f63d3f5e95be685171fa8",
							"name":  "test-job-min-130045",
							"resources": map[string]interface{}{
								"cpu":    "0.3",
								"memory": "768Mi",
							},
							"volumeMounts": []interface{}{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepMerge(tt.newData, tt.oldData)

			// Compare results
			if !mapsEqual(result, tt.expected) {
				t.Errorf("DeepMerge() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// mapsEqual compares two maps for equality
func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v1 := range a {
		v2, ok := b[k]
		if !ok {
			return false
		}

		if !valuesEqual(v1, v2) {
			return false
		}
	}

	return true
}

// valuesEqual compares two values for equality
func valuesEqual(a, b interface{}) bool {
	// Handle nil values
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle maps
	aMap, aIsMap := a.(map[string]interface{})
	bMap, bIsMap := b.(map[string]interface{})
	if aIsMap && bIsMap {
		return mapsEqual(aMap, bMap)
	}

	// Handle slices
	aSlice, aIsSlice := a.([]interface{})
	bSlice, bIsSlice := b.([]interface{})
	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if !valuesEqual(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}

	// Handle primitive types
	return a == b
}
