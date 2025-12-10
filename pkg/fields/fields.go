package fields

import (
	"fmt"
	"reflect"
)

func ValidateFields(data interface{}, model interface{}) error {
	dataType := reflect.TypeOf(data)
	modelType := reflect.TypeOf(model)
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if dataType.Kind() != reflect.Struct || modelType.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs")
	}
	modelFields := make(map[string]bool)
	for i := 0; i < modelType.NumField(); i++ {
		jsonTag := modelType.Field(i).Tag.Get("json")
		if jsonTag != "" {
			modelFields[jsonTag] = true
		}
		modelFields[modelType.Field(i).Name] = true
	}
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			if !modelFields[jsonTag] {
				return fmt.Errorf("field '%s' does not exist in model", jsonTag)
			}
		}
		if !modelFields[field.Name] {
			return fmt.Errorf("field '%s' does not exist in model", field.Name)
		}
	}
	return nil
}
func GetOrDefault(value string, def string) string {
	if value == "" {
		return def
	}
	return value
}
