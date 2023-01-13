package common

import (
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
)

func CopyAny(from, to interface{}) {
	t := reflect.ValueOf(to).Elem()
	if !t.CanSet() {
		return
	}
	f := reflect.ValueOf(from)
	if isBaseNil(f) {
		return
	}
	copyAny(f, t)
}

func setBaseValue(from, to reflect.Value) {
	if isBaseNil(from) {
		return
	}
	var l int
	t := to.Type()
	for t.Kind() == reflect.Ptr {
		l++
		t = t.Elem()
	}
	v := baseValue(from)
	for i := 0; i < l; i++ {
		t := reflect.New(v.Type())
		t.Elem().Set(v)
		v = t
	}
	to.Set(v)
}

func copyAny(from, to reflect.Value) {
	if !to.CanSet() {
		return
	}
	if isBaseNil(from) {
		return
	}
	btfrom := baseType(from.Type())
	btto := baseType(to.Type())
	if typeEq(btfrom, btto) {
		setBaseValue(from, to)
		return
	}
	if _, ok := wrapType[btto.String()]; ok { // string -> wrapperspb.StringValue
		val := reflect.New(btto).Elem()
		copyAny(from, val.FieldByName("Value"))
		setBaseValue(val, to)
		return
	}
	if _, ok := wrapType[btfrom.String()]; ok { // wrapperspb.StringValue -> string
		copyAny(baseValue(from).FieldByName("Value"), to)
		return
	}
	if btfrom.Kind() == reflect.Struct && btto.Kind() == reflect.Struct {
		copyStruct(baseValue(from), baseValue(to))
		return
	}
	if btfrom.Kind() == reflect.Slice && btto.Kind() == reflect.Slice {
		copySlice(baseValue(from), baseValue(to))
		return
	}
	if btto.Kind() == reflect.String {
		if isBaseNil(to) {
			to.Set(getBaseZeroValue(baseType(to.Type())))
		}
		setBaseValue(reflect.ValueOf(toString(from)), to)
		return
	}
	if toNumber(from, to) {
		return
	}
}

func getBaseZeroValue(t reflect.Type) reflect.Value {
	var l int
	for t.Kind() == reflect.Ptr {
		l++
		t = t.Elem()
	}
	v := reflect.Zero(t)
	for i := 0; i < l; i++ {
		t := reflect.New(v.Type())
		t.Elem().Set(v)
		v = t
	}
	return v
}

func isBaseNil(v reflect.Value) bool {
	for {
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		case reflect.Invalid:
			return true
		default:
			return isNil(v)
		}
	}
}

func baseType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func typeEq(t1, t2 reflect.Type) bool {
	return t1.String() == t2.String()
}

func isNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	}
	return false
}

func baseValue(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func copyStruct(from, to reflect.Value) {
	toType := to.Type()
	fromType := from.Type()
	n := to.NumField()
	for i := 0; i < n; i++ {
		toFieldType := toType.Field(i)
		if _, found := fromType.FieldByName(toFieldType.Name); !found {
			continue
		}
		copyAny(from.FieldByName(toFieldType.Name), to.Field(i))
	}
}

func copySlice(from, to reflect.Value) {
	size := from.Len()
	temp := reflect.MakeSlice(to.Type(), 0, size)
	elemTo := to.Type().Elem()
	for i := 0; i < size; i++ {
		itemTo := getBaseZeroValue(elemTo)
		copyAny(from.Index(i), itemTo)
		temp = reflect.Append(temp, itemTo)
	}
	to.Set(temp)
}

func toString(value reflect.Value) string {
	if value.Kind() == reflect.Slice {
		switch value.Type().String() {
		case "[]uint8":
			return string(value.Interface().([]uint8))
		case "[]int32":
			return string(value.Interface().([]int32))
		}
	}
	return fmt.Sprint(value.Interface())
}

func toNumber(from1, to1 reflect.Value) bool {
	if isBaseNil(to1) {
		to1.Set(getBaseZeroValue(to1.Type()))
	}
	from := baseValue(from1)
	to := baseValue(to1)
	switch from.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch to.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			to.SetInt(from.Int())
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			to.SetUint(uint64(from.Int()))
			return true
		case reflect.Float64, reflect.Float32:
			to.SetFloat(float64(from.Int()))
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch to.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			to.SetInt(int64(from.Uint()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			to.SetInt(int64(from.Uint()))
			return true
		case reflect.Float64, reflect.Float32:
			to.SetFloat(float64(from.Uint()))
			return true
		}
	case reflect.Float64, reflect.Float32:
		switch to.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			to.SetInt(int64(from.Float()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			to.SetUint(uint64(from.Float()))
			return true
		case reflect.Float64, reflect.Float32:
			to.SetFloat(from.Float())
			return true
		}
	}
	return false
}

func typeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

var wrapType = map[string]struct{}{
	typeName(wrapperspb.DoubleValue{}): {},
	typeName(wrapperspb.FloatValue{}):  {},
	typeName(wrapperspb.Int64Value{}):  {},
	typeName(wrapperspb.UInt64Value{}): {},
	typeName(wrapperspb.Int32Value{}):  {},
	typeName(wrapperspb.UInt32Value{}): {},
	typeName(wrapperspb.BoolValue{}):   {},
	typeName(wrapperspb.StringValue{}): {},
	typeName(wrapperspb.BytesValue{}):  {},
}
