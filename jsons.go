package flyers

import "encoding/json"

func ToJson(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func FromJson(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

func FromJsonMap(s string) (map[string]any, error) {
	m := make(map[string]any)
	err := FromJson(s, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
