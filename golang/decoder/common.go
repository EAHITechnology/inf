package decoder

import (
	"fmt"
	"reflect"
)

func setField(obj interface{}, m map[string]interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		fieldInfo := structValue.Type().Field(i)
		tag := fieldInfo.Tag
		label := tag.Get("json")
		if label == "" {
			label = fieldInfo.Name
		}

		if _, ok := m[label]; !ok {
			continue
		}
		structFieldValue := structValue.Field(i)
		if !structFieldValue.IsValid() {
			return nil
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", label)
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(m[label])
		if val.Kind() == reflect.Ptr || val.Kind() == reflect.Map {
			continue
		}
		if structFieldType != val.Type() {
			return fmt.Errorf("Provided value type didn't match obj field type")
		}
		structFieldValue.Set(val)
	}
	return nil
}

/*
It is a structure pointer waiting to be converted.
The obj need set json label,eg:
type Test struct{
	Name string `json:"name"`
	Age  int    `json:"age"`
}
*/
func FillStruct(m map[string]interface{}, obj interface{}) error {
	if m == nil {
		return nil
	}
	if err := setField(obj, m); err != nil {
		return err
	}
	return nil
}
