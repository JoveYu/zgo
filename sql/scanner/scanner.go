package scanner

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

var (
	StructTag string = "zdb"
)

func ScanStruct(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}

	v = v.Elem()
	tp := v.Type()

	switch tp.Kind() {
	case reflect.Slice:
		for rows.Next() {
			obj := reflect.New(tp.Elem())

			err := scanOne(rows, obj.Interface())
			if err != nil {
				return err
			}

			v.Set(reflect.Append(v, obj.Elem()))
		}

	case reflect.Struct:
		if rows.Next() {
			err := scanOne(rows, dest)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unknow dest")
	}
	return nil
}

func scanOne(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	v = v.Elem()
	tp := v.Type()

	if v.Kind() != reflect.Struct {
		return errors.New("dest is not struct")
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	fields := make([]interface{}, len(cols))

	for idx, col := range cols {
		ok := false
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			tag := tp.Field(i).Tag.Get(StructTag)
			if tag == col {
				ok = true
				fields[idx] = f.Addr().Interface()
				break
			}
		}
		if !ok {
			return errors.New(fmt.Sprintf("not find field [%s]", col))
		}
	}

	err = rows.Scan(fields...)
	if err != nil {
		return err
	}

	return nil
}
