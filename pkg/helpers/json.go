package helpers

import (
	"encoding/json"
	"os"
)

func ReadJSONFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func WriteJSONFile(filename string, v interface{}, intend ...bool) error {
	var (
		data []byte
		err  error
	)
	if len(intend) > 0 && intend[0] {
		data, err = json.MarshalIndent(v, "", "\t")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
