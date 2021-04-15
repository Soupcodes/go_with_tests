package main

import (
	"reflect"
)

func Walk(x interface{}, fn func(string)) {
	val := getValue(x)

	numValues := 0
	var getField func(int) reflect.Value

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		numValues = val.NumField() // Get all the fields in the struct
		getField = val.Field       // assign this function to a variable
	case reflect.Slice, reflect.Array:
		numValues = val.Len()
		getField = val.Index
	}

	for i := 0; i < numValues; i++ { // Loops through all values in struct / slice
		Walk(getField(i).Interface(), fn) // Recursively get the value and try to run fn
	}
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
