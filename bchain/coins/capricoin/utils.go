package capricoin

import (
	"bytes"
	"encoding/json"
)

func removeDuplicateJSONKeys(inJSON []byte) ([]byte, error) {
	val, err := decode(json.NewDecoder(bytes.NewReader(inJSON)))
	if err != nil {
		return nil, err
	}

	outJSON, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	return outJSON, nil
}

// Ignores any duplicate key/value pairs in JSON objects, keeps the first occurrences only
func decode(d *json.Decoder) (interface{}, error) {
	// Get next token from JSON
	t, err := d.Token()
	if err != nil {
		return nil, err
	}

	delim, ok := t.(json.Delim)

	// Return simple arr (strings, numbers, bool, nil)
	if !ok {
		return t, nil
	}

	switch delim {
	case '{':
		dict := make(map[string]interface{})
		for d.More() {
			// Get field key
			t, err := d.Token()
			if err != nil {
				return nil, err
			}
			key := t.(string)

			value, err := decode(d)
			if err != nil {
				return nil, err
			}
			// Ignore duplicate keys, store only the first occurence
			if _, ok := dict[key]; !ok {
				dict[key] = value
			}
		}
		// Consume trailing }
		if _, err := d.Token(); err != nil {
			return nil, err
		}
		return dict, nil

	case '[':
		var arr []interface{}
		i := 0
		for d.More() {
			value, err := decode(d)
			if err != nil {
				return nil, err
			}
			i++
			arr = append(arr, value)
		}
		// Consume trailing ]
		if _, err := d.Token(); err != nil {
			return nil, err
		}
		return arr, nil
	}
	return nil, nil
}
