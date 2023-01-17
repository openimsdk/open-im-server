package common

import (
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
)

func CopyAny(from, to interface{}) {
	copyAny(reflect.ValueOf(from), reflect.Indirect(reflect.ValueOf(to)))
}

func copyAny(from, to reflect.Value) {
	if !to.CanSet() {
		return
	}
	if isBaseNil(from) {
		return
	}
	if isBaseNil(to) {
		to.Set(getBaseZeroValue(to.Type()))
	}
	btFrom := baseType(from.Type())
	btTo := baseType(to.Type())
	if btTo.Kind() == reflect.Interface || typeEq(btFrom, btTo) {
		setBaseValue(from, to)
		return
	}
	if _, ok := wrapType[btTo.String()]; ok { // string -> wrapperspb.StringValue
		val := reflect.New(btTo).Elem()
		copyAny(from, val.FieldByName("Value"))
		setBaseValue(val, to)
		return
	}
	if _, ok := wrapType[btFrom.String()]; ok { // wrapperspb.StringValue -> string
		copyAny(baseValue(from).FieldByName("Value"), to)
		return
	}
	if btFrom.Kind() == reflect.Struct && btTo.Kind() == reflect.Struct {
		copyStruct(baseValue(from), baseValue(to))
		return
	}
	if btFrom.Kind() == reflect.Slice && btTo.Kind() == reflect.Slice {
		copySlice(baseValue(from), baseValue(to))
		return
	}
	if btFrom.Kind() == reflect.Map && btTo.Kind() == reflect.Map {
		copyMap(baseValue(from), baseValue(to))
		return
	}
	if btTo.Kind() == reflect.String {
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

func getBaseZeroValue(t reflect.Type) reflect.Value {
	var l int
	for t.Kind() == reflect.Ptr {
		l++
		t = t.Elem()
	}
	v := reflect.New(t)
	if l == 0 {
		v = v.Elem()
	} else {
		for i := 1; i < l; i++ {
			t := reflect.New(v.Type())
			t.Elem().Set(v)
			v = t
		}
	}
	r := reflect.New(v.Type()).Elem()
	r.Set(v)
	return r
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
		var itemTo reflect.Value
		if item := from.Index(i); isBaseNil(item) {
			itemTo = reflect.Zero(elemTo)
		} else {
			itemTo = getBaseZeroValue(elemTo)
			copyAny(from.Index(i), itemTo)
		}
		temp = reflect.Append(temp, itemTo)
	}
	to.Set(temp)
}

func copyMap(from, to reflect.Value) {
	to.Set(reflect.MakeMap(to.Type()))
	toTypeKey := to.Type().Key()
	toTypeVal := to.Type().Elem()
	for r := from.MapRange(); r.Next(); {
		key := getBaseZeroValue(toTypeKey)
		copyAny(r.Key(), key)
		var val reflect.Value
		fVal := r.Value()
		if isBaseNil(fVal) {
			val = reflect.Zero(toTypeVal)
		} else {
			val = getBaseZeroValue(toTypeVal)
			copyAny(fVal, val)
		}
		to.SetMapIndex(key, val)
	}
}

func toString(value reflect.Value) string {
	if value.Kind() == reflect.Slice {
		switch value.Type().String() {
		case "[]uint8": // []byte -> []uint8
			return string(value.Interface().([]uint8))
		case "[]int32": // []rune -> []int32
			return string(value.Interface().([]int32))
		}
	}
	return fmt.Sprint(value.Interface())
}

func toNumber(from, to reflect.Value) bool {
	initTo := func() {
		if isBaseNil(to) {
			to.Set(getBaseZeroValue(to.Type()))
		}
	}
	switch baseValue(from).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch baseValue(to).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			initTo()
			baseValue(to).SetInt(baseValue(from).Int())
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			initTo()
			baseValue(to).SetUint(uint64(baseValue(from).Int()))
			return true
		case reflect.Float64, reflect.Float32:
			initTo()
			baseValue(to).SetFloat(float64(baseValue(from).Int()))
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch baseValue(to).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			initTo()
			baseValue(to).SetInt(int64(baseValue(from).Uint()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			initTo()
			baseValue(to).SetInt(int64(baseValue(from).Uint()))
			return true
		case reflect.Float64, reflect.Float32:
			initTo()
			baseValue(to).SetFloat(float64(baseValue(from).Uint()))
			return true
		}
	case reflect.Float64, reflect.Float32:
		switch baseValue(to).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			initTo()
			baseValue(to).SetInt(int64(baseValue(from).Float()))
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			initTo()
			baseValue(to).SetUint(uint64(baseValue(from).Float()))
			return true
		case reflect.Float64, reflect.Float32:
			initTo()
			baseValue(to).SetFloat(baseValue(from).Float())
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
