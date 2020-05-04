package enums

import (
	"fmt"
	"reflect"
)

type enumListMap map[interface{}]interface{}

// Enum represents a fixed set of values of the same type that are valid for a given
// attribute.
type Enum struct {
	enumList      enumListMap
	valueType     reflect.Type
	convertValues bool
}

func isHashable(t reflect.Type) bool {
	kind := t.Kind()
	for _, k := range []reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
	} {
		if kind == k {
			return true
		}
	}

	return false
}

// New creates a new Enum comprised of provided values. It's Immutable.
func New(values ...interface{}) Enum {
	e := Enum{
		enumList:      enumListMap{},
		convertValues: false,
	}

	if len(values) == 0 {
		panic("Enum is empty.")
	}

	for _, value := range values {
		t := reflect.TypeOf(value)
		if e.valueType == nil {
			if !isHashable(t) {
				panic("Enum values must be of hashable type.")
			}
			e.valueType = t
		} else if t != e.valueType {
			panic(fmt.Sprintf("Found values of different kinds in enum: %v and %v.",
				e.valueType, t))
		}
		if _, found := e.enumList[value]; found {
			panic(fmt.Sprintf("Found repeated value on enum: %v.", value))
		}
		e.enumList[value] = true
	}

	return e
}

// NewConvert creates a new Enum comprised of provided values. convertValues will be set
// to true, meaning that it'll allow checking for validity on any value whose type is
// convertible to this Enum own type. It's Immutable.
func NewConvert(values ...interface{}) Enum {
	e := New(values...)
	e.convertValues = true
	return e
}

func (e Enum) isSupportedValueType(t reflect.Type) bool {
	return t == e.valueType || (e.convertValues && t.ConvertibleTo(e.valueType))
}

func (e Enum) supportedValueTypeOrPanic(value interface{}) {
	t := reflect.TypeOf(value)
	if !e.isSupportedValueType(t) {
		convertValuesClause := ""
		if e.convertValues {
			convertValuesClause = " or convertible to it"
		}
		panic(fmt.Sprintf("Invalid type: %v. Should be %v%s.", t, e.valueType,
			convertValuesClause))
	}
}

// IsValid returns true if provided value exists within the Enum.
func (e Enum) IsValid(value interface{}) bool {
	e.supportedValueTypeOrPanic(value)

	if e.convertValues && reflect.TypeOf(value) != e.valueType {
		value = reflect.ValueOf(value).Convert(e.valueType).Interface()
	}

	if _, found := e.enumList[value]; !found {
		return false
	}

	return true
}

// IsAnyValid returns true if any of the provided values exists within the Enum.
func (e Enum) IsAnyValid(values ...interface{}) bool {
	for _, value := range values {
		if valid := e.IsValid(value); valid {
			return true
		}
	}
	return false
}

// AreAllValid returns true if all of the provided values exist within the Enum.
func (e Enum) AreAllValid(values ...interface{}) bool {
	for _, value := range values {
		if valid := e.IsValid(value); !valid {
			return false
		}
	}
	return true
}

func (e Enum) supportedFunctionTypeOrPanic(targetFunc reflect.Value, variadic bool) {
	inType := e.valueType
	if variadic {
		inType = reflect.SliceOf(inType)
	}
	expectedFuncType := reflect.FuncOf(
		[]reflect.Type{inType},
		[]reflect.Type{reflect.TypeOf(true)},
		variadic)
	convertValuesClause := ""
	if e.convertValues {
		convertValuesClause = fmt.Sprintf(
			"\nArg input type can be %v, or type convertible to it.", e.valueType)
	}

	ptrType := targetFunc.Type()

	panicWithDetails := func(details string) {
		panic(fmt.Sprintf(
			"targetFunc should be a Ptr to a %v function. Received %v - %s%s",
			expectedFuncType, ptrType, details, convertValuesClause))
	}

	if ptrType.Kind() != reflect.Ptr {
		panicWithDetails("not a pointer")
	}

	t := targetFunc.Elem().Type()

	if t.Kind() != reflect.Func {
		panicWithDetails("not a function")
	}

	if t.NumIn() != 1 ||
		(!variadic && !e.isSupportedValueType(t.In(0))) ||
		(variadic && (t.In(0).Kind() != reflect.Slice ||
			!e.isSupportedValueType(t.In(0).Elem()))) {
		panicWithDetails("check IN args")
	}

	if t.NumOut() != 1 || t.Out(0) != expectedFuncType.Out(0) {
		panicWithDetails("check OUT args")
	}
}

func (e Enum) setTypedMethod(targetFunc interface{}, method reflect.Value) {
	v := reflect.ValueOf(targetFunc)
	variadic := method.Type().IsVariadic()
	e.supportedFunctionTypeOrPanic(v, variadic)
	v.Elem().Set(reflect.MakeFunc(v.Elem().Type(), func(args []reflect.Value) []reflect.Value {
		if !variadic {
			return method.Call(args)
		}

		a := []interface{}{}
		for i := 0; i < args[0].Len(); i++ {
			a = append(a, args[0].Index(i).Interface())
		}
		args[0] = reflect.ValueOf(a)
		return method.CallSlice(args)
	}))
}

// SetTypedIsValid sets the provided function to be a typed version of this enum.IsValid.
func (e Enum) SetTypedIsValid(targetFunc interface{}) {
	e.setTypedMethod(targetFunc, reflect.ValueOf(e.IsValid))
}

// SetTypedIsAnyValid sets the provided function to be a typed version of this
// enum.IsAnyValid. *WARNING* this function will need to internally cast each provided
// values to interface{} and store them in an aux array, adding to time and spatial
// complexity.
func (e Enum) SetTypedIsAnyValid(targetFunc interface{}) {
	e.setTypedMethod(targetFunc, reflect.ValueOf(e.IsAnyValid))
}

// SetTypedAreAllValid sets the provided function to be a typed version of this
// enum.AreAllValid. *WARNING* this function will need to internally cast each provided
// values to interface{} and store them in an aux array, adding to time and spatial
// complexity.
func (e Enum) SetTypedAreAllValid(targetFunc interface{}) {
	e.setTypedMethod(targetFunc, reflect.ValueOf(e.AreAllValid))
}
