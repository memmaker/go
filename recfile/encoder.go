package recfile

import (
	"errors"
	"io"
	"reflect"
	"sync"
)

type Encoder struct {
	writer io.Writer
	mutex  sync.Mutex
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{writer: writer}
}

func (enc *Encoder) Encode(e any) error {
	return enc.encodeValue(reflect.ValueOf(e))
}

func (enc *Encoder) encodeValue(value reflect.Value) error {
	if value.Kind() == reflect.Invalid {
		return errors.New("rec: cannot encode nil value")
	}
	if value.Kind() == reflect.Pointer && value.IsNil() {
		panic("rec: cannot encode nil pointer of type " + value.Type().String())
	}

	// Make sure we're single-threaded through here, so multiple
	// goroutines can share an encoder.
	enc.mutex.Lock()
	defer enc.mutex.Unlock()

	typeName := value.Type().Name()

	record, err := serialize("", value)
	if err != nil {
		return err
	}
	return WriteMulti(enc.writer, map[string][]Record{typeName: {record}})
}

func ValueToRecord(value any) (Record, error) {
	return serialize("", reflect.ValueOf(value))
}

func serialize(name string, v reflect.Value) (Record, error) {
	var record Record
	var err error
	switch v.Kind() {
	case reflect.String:
		fieldName := name
		record = append(record, Field{Name: fieldName, Value: v.String()})
	case reflect.Int:
		fieldName := name
		record = append(record, Field{Name: fieldName, Value: Int64Str(v.Int())})
	case reflect.Float64:
		fieldName := name
		record = append(record, Field{Name: fieldName, Value: FloatStr(v.Float())})
	case reflect.Bool:
		fieldName := name
		record = append(record, Field{Name: fieldName, Value: BoolStr(v.Bool())})
	case reflect.Ptr:
		if v.IsNil() {
			err = errors.New("rec: cannot encode nil pointer of type " + v.Type().String())
		} else {
			record, err = serialize("", v.Elem())
		}
	case reflect.Struct:
		record, err = serializeStruct(name, v)
	}
	return record, err
}

func serializeStruct(name string, v reflect.Value) (Record, error) {
	t := v.Type()
	var record Record
	for i := 0; i < t.NumField(); i++ {
		structField := v.Type().Field(i)
		fieldName := structField.Name
		if name != "" {
			fieldName = name + "_" + fieldName
		}
		valueFields, err := serialize(fieldName, v.Field(i))
		if err != nil {
			return nil, err
		}
		record = append(record, valueFields...)
	}
	return record, nil
}
