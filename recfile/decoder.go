package recfile

/* Out of order until further notice
func deserialize(name string, rec Record, t reflect.Type) (interface{}, error) {
	var err error
	var value interface{}
	switch t.Kind() {
	case reflect.Ptr:
		p := reflect.New(t.Elem())
		de, deErr := deserialize(name, rec, t.Elem())
		if deErr != nil {
			return nil, deErr
		}
		p.Elem().Set(reflect.ValueOf(de))
		value, err = p.Interface(), nil
	case reflect.Struct:
		value, err = deserializeStruct(name, rec, t)
	}
	return value, err
}

func deserializeStruct(name string, rec Record, t reflect.Type) (interface{}, error) {
	v := reflect.New(t).Elem()

	for _, field := range rec {
		fieldName := field.Name
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv, err := deserialize(name, rec, t.Field(i).Type)
		if err != nil {
			return nil, err
		}
		v.Field(i).Set(reflect.ValueOf(fv).Convert(field.Type))
	}
	return v.Interface(), nil
}
*/
