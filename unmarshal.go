package queries

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(dst interface{}, src url.Values) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("destination interface was a nil pointer or not a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("destination interface did not point to a struct")
	}
	return unmarshalStruct(rv, src)
}

func unmarshalStruct(dst reflect.Value, src url.Values) error {
	dstTy := dst.Type()
	for i := 0; i < dstTy.NumField(); i++ {
		field := dstTy.Field(i)
		tag := field.Tag.Get("json")
		required := !strings.HasSuffix(tag, ",omitempty")

		if off := strings.IndexByte(tag, ','); off != -1 {
			tag = tag[:off]
		}

		if tag == "" {
			return fmt.Errorf("field was missing a tag")
		}

		src := src[tag]
		dst := dst.Field(i)

		if len(src) == 0 {
			if required {
				return fmt.Errorf("required field was missing")
			}
			continue
		}

		if err := unmarshalValue(dst, src[0]); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalValue(dst reflect.Value, src string) error {
	if dst.Kind() != reflect.Ptr {
		return unmarshalScalarValue(dst, src)
	}
	v := reflect.New(dst.Type().Elem())
	err := unmarshalScalarValue(v.Elem(), src)
	if err == nil {
		dst.Set(v)
	}
	return err
}

func unmarshalScalarValue(dst reflect.Value, src string) error {
	switch dst.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(src, 10, dst.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse int")
		}
		dst.SetInt(v)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(src, 10, dst.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse int")
		}
		dst.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(src, dst.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse int")
		}
		dst.SetFloat(v)
	case reflect.String:
		dst.SetString(src)
	}
	return nil
}
