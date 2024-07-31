package visible

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

// delete all fields that doesn't have the visible tag value
//
// example:
//
//	type Data struct {
//		ID    int
//		Name  string `json:"name"`
//		Data1 string `json:"data1" visible:"alice,bob"`
//		Data2 string `json:"data2" visible:"bob"`
//		HiddenData string `json:"-"`
//	}
//
// CleanStruct(data, "alice") -> {id: 1, name: "name", data1: "data1"}	// only alice_data is visible
//
// CleanStruct(data, "bob") -> {id: 1, name: "name", data1: "data1", data2: "data2"}	// only bob_data is visible
//
// CleanStruct(data, "charlie") -> {id: 1, name: "name"}	// no data is visible
func CleanStruct(data interface{}, visiblityValidator string) (map[string]interface{}, error) {

	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if !structs.IsStruct(data) {
		return nil, fmt.Errorf("data is not a struct")
	}

	fields := structs.Fields(data)
	res := make(map[string]interface{})
	for _, field := range fields {

		if isJsonIgnored(field) {
			continue
		}

		tag := extractVisibleTag(field)
		if tag == "" {
			res[field.Name()] = field.Value()
			continue
		}

		if isFieldVisible(tag, visiblityValidator) {
			res[field.Name()] = field.Value()
			continue
		}

	}
	return res, nil
}

func isJsonIgnored(field *structs.Field) bool {
	return field.Tag("json") == "-"
}

func extractVisibleTag(field *structs.Field) string {
	return field.Tag("visible")
}

func isFieldVisible(tag string, validator string) bool {
	tags := strings.Split(tag, ",")

	for _, t := range tags {
		if t == validator {
			return true
		}
	}
	return false
}
