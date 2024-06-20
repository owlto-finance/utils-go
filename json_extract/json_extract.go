package json_extract

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractValueFromJSON(jsonData string, path string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	// 分割路径 "estimation.costsDetails[0].chain" -> ["estimation", "costsDetails[0]", "chain"]
	pathParts := strings.Split(path, ".")

	current := data
	for _, part := range pathParts {
		if strings.Contains(part, "[") {
			// 处理数组部分
			re := regexp.MustCompile(`(\w+)\[(\d+)\]`)
			matches := re.FindStringSubmatch(part)
			if len(matches) != 3 {
				return nil, fmt.Errorf("invalid path segment: %s", part)
			}

			key := matches[1]
			index, _ := strconv.Atoi(matches[2])

			// 转到指定的key
			if mapped, ok := current.(map[string]interface{}); ok {
				if value, found := mapped[key]; found {
					// 转到数组的指定索引
					if sliced, ok := value.([]interface{}); ok {
						if index < len(sliced) {
							current = sliced[index]
						} else {
							return nil, fmt.Errorf("index out of range: %d", index)
						}
					} else {
						return nil, fmt.Errorf("not an array: %s", key)
					}
				} else {
					return nil, fmt.Errorf("key not found: %s", key)
				}
			} else {
				return nil, fmt.Errorf("not a map before key: %s", key)
			}
		} else {
			if mapped, ok := current.(map[string]interface{}); ok {
				if value, found := mapped[part]; found {
					current = value
				} else {
					return nil, fmt.Errorf("key not found: %s", part)
				}
			} else {
				return nil, fmt.Errorf("not a map at key: %s", part)
			}
		}
	}
	return current, nil
}

func ExtractStringValueFromJSON(jsonData string, path string) (string, error) {
	value, err := ExtractValueFromJSON(jsonData, path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func ExtractStringValueFromObj(obj interface{}, path string) (string, error) {
	jsonData, _ := json.Marshal(obj)
	value, err := ExtractValueFromJSON(string(jsonData), path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func ExtractInt64ValueFromJSON(jsonData string, path string) (int64, error) {
	value, err := ExtractValueFromJSON(jsonData, path)
	if err != nil {
		return 0, err
	}

	if floatValue, ok := value.(float64); ok {
		return int64(floatValue), nil
	} else {
		return 0, fmt.Errorf("not an integer: %s", path)
	}
}

func ExtractInt64ValueFromObj(obj interface{}, path string) (int64, error) {
	jsonData, _ := json.Marshal(obj)
	value, err := ExtractValueFromJSON(string(jsonData), path)
	if err != nil {
		return 0, err
	}

	if floatValue, ok := value.(float64); ok {
		return int64(floatValue), nil
	} else {
		return 0, fmt.Errorf("not an integer: %s", path)
	}
}

func ExtractInt32ValueFromJSON(jsonData string, path string) (int32, error) {
	value, err := ExtractValueFromJSON(jsonData, path)
	if err != nil {
		return 0, err
	}

	if floatValue, ok := value.(float64); ok {
		return int32(floatValue), nil
	} else {
		return 0, fmt.Errorf("not an integer: %s", path)
	}
}

func ExtractInt32ValueFromObj(obj interface{}, path string) (int32, error) {
	jsonData, _ := json.Marshal(obj)
	value, err := ExtractValueFromJSON(string(jsonData), path)
	if err != nil {
		return 0, err
	}

	if floatValue, ok := value.(float64); ok {
		return int32(floatValue), nil
	} else {
		return 0, fmt.Errorf("not an integer: %s", path)
	}
}

func ExtractSliceValueFromJSON(jsonData string, path string) ([]interface{}, error) {
	value, err := ExtractValueFromJSON(jsonData, path)
	if err != nil {
		return nil, err
	}
	if slice, ok := value.([]interface{}); ok {
		return slice, nil
	} else {
		return nil, fmt.Errorf("not an array: %s", path)
	}
}

func ExtractSliceValueFromObj(obj interface{}, path string) ([]interface{}, error) {
	jsonData, _ := json.Marshal(obj)
	value, err := ExtractValueFromJSON(string(jsonData), path)
	if err != nil {
		return nil, err
	}
	if slice, ok := value.([]interface{}); ok {
		return slice, nil
	} else {
		return nil, fmt.Errorf("not an array: %s", path)
	}
}
