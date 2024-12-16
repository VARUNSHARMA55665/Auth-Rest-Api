package helpers

import (
	"encoding/json"
	"log"
)

func LogStructAsJSON(data interface{}) interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("LogStructAsJSON Error marshaling JSON: ", err)
		// if there is error in marshalling then atleast return the exact same packet
		return data
	}

	return string(jsonData)
}
