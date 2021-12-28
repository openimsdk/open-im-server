package utils

import (
	"fmt"
	"github.com/jinzhu/copier"
	"reflect"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)

	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return err
	}
	av = reflect.ValueOf(av.Interface())

	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		err = fmt.Errorf("no fields to copy")
		return err
	}

	for i := 0; i < len(_fields); i++ {
		name := _fields[i]

		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		}
	}
	return nil
}
