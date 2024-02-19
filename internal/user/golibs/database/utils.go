package database

import (
	"reflect"
	"strconv"
	"strings"
)

func FieldMap(e interface{}) ([]string, []interface{}) {
	var fieldNames []string
	var fieldValues []interface{}
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Tag.Get("db")
		fieldValue := v.Field(i).Addr().Interface()
		fieldNames = append(fieldNames, fieldName)
		fieldValues = append(fieldValues, fieldValue)
	}

	return fieldNames, fieldValues
}

func GeneratePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	var builder strings.Builder
	sep := ", "
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		builder.WriteString("$" + strconv.Itoa(i) + sep)
	}

	return builder.String()
}
