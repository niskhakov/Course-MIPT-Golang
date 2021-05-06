package params

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	DEBUG = false
)

var (
	ErrNotPointer = errors.New("object not a pointer")
	ErrNotSupportedFieldType = errors.New("field type is not supported")
)

var (
	typeIntSlice []int
	typeIntPtr *int
	typeStrSlice []string
	typeStrPtr *string
	typeBoolSlice []bool
	typeBoolPtr	*bool
	typeInt int
	typeStr string
	typeBool bool
)

var dbg *log.Logger

func init() {
	var out io.Writer 
	if DEBUG {
		out = os.Stdout
	} else {
		out = ioutil.Discard
	}
	dbg = log.New(out, "DEBUG: ", log.Lmicroseconds)
}


func Unpack(values url.Values, to interface{}) error {

	tagNameMap := make(map[string]int)

	t := reflect.TypeOf(to)
	if t.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	v := reflect.ValueOf(to).Elem()
	t = t.Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
	    first := field.Name[:1] 
		if strings.ToUpper(first) != first {
			dbg.Printf("Private field encountered - %s, ignoring\n", field.Name)
			continue
		}
		tag := field.Tag.Get("http")
		if len(tag) ==  0 {
			tag = strings.ToLower(field.Name)
		}
		tagNameMap[tag] = i
		dbg.Printf("%d: %s -> %s\n", tagNameMap[tag], field.Name, tag)
	}

	for key, val := range values {
		if indx, ok := tagNameMap[key]; ok {
			dbg.Println(key, indx, val) 
			typeObj := v.Type().Field(indx).Type
			switch typeObj {
				case reflect.TypeOf(typeBool), reflect.TypeOf(typeStr), reflect.TypeOf(typeInt):
					setReflectPrimitive(typeObj, v.Field(indx), values[key])	

				case reflect.TypeOf(typeIntPtr), reflect.TypeOf(typeStrPtr), reflect.TypeOf(typeBoolPtr):
					setReflectPtr(typeObj, v.Field(indx), values[key])

				case reflect.TypeOf(typeIntSlice): 
					intSliceReflect := reflect.MakeSlice(reflect.TypeOf(typeIntSlice), len(values[key]), len(values[key]))
					v.Field(indx).Set(intSliceReflect)
					
					slice := intSliceReflect.Interface().([]int)
					for i, v := range values[key] {
						intv, err := strconv.Atoi(v)
						if err != nil {
							continue
						}
						slice[i] = intv
					}
					dbg.Printf("int slice set to: %v\n ", slice)
				case reflect.TypeOf(typeStrSlice):
					strSliceReflect := reflect.MakeSlice(reflect.TypeOf(typeStrSlice), len(values[key]), len(values[key]))
					v.Field(indx).Set(strSliceReflect)
					
					slice := strSliceReflect.Interface().([]string)
					for i, v := range values[key] {
						slice[i] = v
					}
					dbg.Printf("str slice set to: %v\n ", slice)

				case reflect.TypeOf(typeBoolSlice):
					boolSliceReflect := reflect.MakeSlice(reflect.TypeOf(typeBoolSlice), len(values[key]), len(values[key]))
					v.Field(indx).Set(boolSliceReflect)
					
					slice := boolSliceReflect.Interface().([]bool)
					for i, v := range values[key] {
						boolv, err := strconv.ParseBool(v)
						if err != nil {
							continue
						}
						slice[i] = boolv
					}
					dbg.Printf("bool slice set to: %v\n ", slice)
				default:
					return ErrNotSupportedFieldType
			}
		}
	}
	return nil
}

func setReflectPrimitive(objType reflect.Type, field reflect.Value, vval []string) {
	switch objType {
		case reflect.TypeOf(typeBool):
			if bl, err := strconv.ParseBool(vval[0]); err == nil {
				field.SetBool(bl)
			}

		case reflect.TypeOf(typeStr):
			field.SetString(vval[0])

		case reflect.TypeOf(typeInt):
			if it, err := strconv.ParseInt(vval[0], 10, 64); err == nil {
				field.SetInt(it)
			}
	}

	dbg.Printf("field %s set to: %v\n", field.Type().Name(), field)
}

func setReflectPtr(objType reflect.Type, field reflect.Value, vval []string) {
	ptr := reflect.New(field.Type().Elem())
	field.Set(ptr)

	switch objType {
		case reflect.TypeOf(typeIntPtr):
			if it, err := strconv.ParseInt(vval[0], 10, 64); err == nil {
				ptr.Elem().SetInt(it)
			}
		case reflect.TypeOf(typeStrPtr):
			ptr.Elem().SetString(vval[0])
		case reflect.TypeOf(typeBoolPtr):
			if bl, err := strconv.ParseBool(vval[0]); err == nil {
				ptr.Elem().SetBool(bl)
			}
	}
	dbg.Printf("some ptr set to: %v\n", field)
}

