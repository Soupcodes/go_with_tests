package main

import (
	"reflect"
)

/*
Traverse arbitrary data types using reflect

The go docs conclude that these are the laws of reflection:
    * Reflection goes from interface value to reflection object.

    * Reflection goes from reflection object to interface value.

    * To modify a reflection object, the value must be settable.

Once you understand these laws reflection in Go becomes much easier to use, although it remains subtle.
It's a powerful tool that should be used with care and avoided unless strictly necessary.
*/

func Walk(x interface{}, fn func(string)) {
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		Walk(value.Interface(), fn)
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		// For value on receive from the channel; continue if ok is true; and after each iteration assign v to the next value on the channel
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			Walk(v.Interface(), fn)
		}
	case reflect.Func:
		result := val.Call(nil)      // Calls the function with the relevant arguments and appends them to a []reflect.Values (result)
		for _, res := range result { // Loop through the []reflect.Values and invoke Walk
			Walk(res.Interface(), fn)
		}
	}
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
