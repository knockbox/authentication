package utils

import "reflect"

// PayloadHasChanges accepts a payload by value (*payload) and checks to see if there
// are any changes present on the interface. Returns true if at least one field is present.
//
// O(n) where n is the number of fields present in the interface
func PayloadHasChanges(payload interface{}) bool {
	val := reflect.ValueOf(payload)

	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsNil() {
			return true
		}
	}

	return false
}
