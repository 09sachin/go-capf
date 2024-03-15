package controllers

import (
	"encoding/json"
	"fmt"
)

func (cs CustomString) MarshalJSON() ([]byte, error) {
	if cs.Valid {
		return json.Marshal(cs.String)
	}
	return json.Marshal("") // Convert null to empty string in JSON
}


func (cs *CustomString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	cs.Valid = true
	cs.String = str
	return nil
}


func formatStringSlice(slice []string) string {
	result := ""
	for i, value := range slice {
		result += fmt.Sprintf("'%s'", value)
		if i < len(slice)-1 {
			result += ", "
		}
	}
	return result
}