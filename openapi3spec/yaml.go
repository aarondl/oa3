package openapi3spec

import (
	"fmt"
	"reflect"
	"strings"
)

// ObjectUnmarshaler is used as a last ditch resort for particularly horrifying
// structs.
type ObjectUnmarshaler interface {
	UnmarshalYAMLObject(intf any) error
}

type mismatchErr struct {
	goType, yamlType string
}

func (m mismatchErr) Error() string {
	return fmt.Sprintf("type mismatch go:(%s) yaml:(%s)", m.goType, m.yamlType)
}

func mismatch(goVal reflect.Value, yamlVal any) error {
	var err mismatchErr
	if goVal.IsValid() {
		err.goType = goVal.Type().Name()
	} else {
		err.goType = "(unknown)"
	}

	if yamlVal == nil {
		err.yamlType = "(nil)"
	} else {
		yamlValue := reflect.ValueOf(yamlVal)
		if yamlValue.IsValid() {
			err.yamlType = yamlValue.Type().Name() + "/" + yamlValue.Type().String()
		} else {
			err.yamlType = "(unknown)"
		}
	}

	return err
}

// UnmarshalYAML completely overrides the typical recursive yaml decoder
// behavior with its own ideas about how to unmarshal in order to handle
// some idiosynchracies in the spec.
func (o *OpenAPI3) UnmarshalYAML(unmarshal func(any) error) error {
	untyped := make(map[string]any)

	if err := unmarshal(untyped); err != nil {
		return err
	}

	return yamlStruct(reflect.ValueOf(o), untyped)
}

func yamlStruct(goValue reflect.Value, yamlObj map[string]any) error {
	if goValue.Kind() == reflect.Ptr {
		goValue = goValue.Elem()
	}
	goType := goValue.Type()

	// This is a reference struct and we special case it by either finding a
	// ref key, or by finding normal values.
	if goType.NumField() == 2 && goType.Field(0).Name == "Ref" && goType.Field(1).Anonymous {
		refIntf, ok := yamlObj["$ref"]
		if ok {
			ref, ok := refIntf.(string)
			if !ok {
				return fmt.Errorf("$ref value in %s was not a string", goType.Name())
			}
			debugln("ref(ref): ", ref)

			goValue.Field(0).SetString(ref)

			return nil
		}

		debugln("ref(emb): ", goValue.Field(1).Type().String())

		return allocAndSet(goValue.Field(1), yamlObj)
	}

	fields := fieldMap("yaml", goType)

	for k, v := range yamlObj {
		field, ok := fields[k]
		if !ok {
			if !strings.HasPrefix(k, "x-") {
				return fmt.Errorf("invalid key for struct (%s): %s", goType.Name(), k)
			}

			extMapField, ok := fields["extensions"]
			if !ok {
				return fmt.Errorf("invalid key for struct (%s; extensions not supported): %s", goType.Name(), k)
			}

			extMapFieldValue := goValue.FieldByIndex(extMapField.Index)
			debugln("extension:", k, reflect.TypeOf(v).String())

			extMapIntf := extMapFieldValue.Interface()
			extMap := extMapIntf.(Extensions)
			if extMap == nil {
				extMap = make(Extensions)
				extMapFieldValue.Set(reflect.ValueOf(extMap))
			}
			extMap[k] = v
			continue
		}

		debugln("structkey:", k, field.Type.String())

		fieldValue := goValue.FieldByIndex(field.Index)

		if err := allocAndSet(fieldValue, v); err != nil {
			return fmt.Errorf("failed to set struct field (%s.%s): %w", goType.Name(), k, err)
		}
	}

	return nil
}

func yamlSlice(goSlice reflect.Value, yamlSlice []any) error {
	for i, v := range yamlSlice {
		goItem := goSlice.Index(i)
		if err := allocAndSet(goItem, v); err != nil {
			return fmt.Errorf("failed to set array index [%d]: %w", i, err)
		}
	}

	return nil
}

func yamlMap(goMap reflect.Value, yamlObject map[string]any) error {
	mapValueType := goMap.Type().Elem()

	for k, v := range yamlObject {
		if mapValueType.Kind() == reflect.Ptr {
			debugln("map(ptr): ", k, goMap.Type().String())
			valType := mapValueType.Elem()

			newVal := reflect.New(valType)
			if err := allocAndSet(newVal.Elem(), v); err != nil {
				return fmt.Errorf("failed setting map key (%s): %w", k, err)
			}
			goMap.SetMapIndex(reflect.ValueOf(k), newVal)
		} else {
			debugln("map(val): ", k, goMap.Type().String())
			newVal := reflect.New(mapValueType)
			if err := allocAndSet(newVal.Elem(), v); err != nil {
				return fmt.Errorf("failed setting map key (%s): %w", k, err)
			}
			goMap.SetMapIndex(reflect.ValueOf(k), newVal.Elem())
		}
	}

	return nil
}

func allocAndSet(val reflect.Value, yamlVal any) error {
	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val.Set(reflect.New(typ))

		// This oddity is to stop the normal yamlStruct from hitting and instead
		// pass any object into the ObjectUnmarshaler after allocation. This
		// is because Go's inherent lack of union struct types.
		goIntf := val.Interface()
		if unmarshaler, ok := goIntf.(ObjectUnmarshaler); ok {
			return unmarshaler.UnmarshalYAMLObject(yamlVal)
		}

		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		yamlArray, ok := yamlVal.([]any)
		if !ok {
			return mismatch(val, yamlVal)
		}

		ln := len(yamlArray)
		if len(yamlArray) == 0 {
			return nil
		}

		slice := reflect.MakeSlice(typ, ln, ln)
		val.Set(slice)

		err := yamlSlice(slice, yamlArray)
		if err != nil {
			return err
		}
	case reflect.Struct:
		yamlObj, ok := yamlVal.(map[string]any)
		if !ok {
			return mismatch(val, yamlVal)
		}

		if err := yamlStruct(val, yamlObj); err != nil {
			return err
		}
	case reflect.Map:
		yamlObj, ok := yamlVal.(map[string]any)
		if !ok {
			return mismatch(val, yamlVal)
		}

		mp := reflect.MakeMapWithSize(typ, len(yamlObj))
		val.Set(mp)

		err := yamlMap(val, yamlObj)
		if err != nil {
			return err
		}
	default:
		if err := yamlPrimitive(val, yamlVal); err != nil {
			return err
		}
	}

	return nil
}

func yamlPrimitive(goPrimitive reflect.Value, yamlPrimitive any) error {
	switch goPrimitive.Kind() {
	case reflect.Float32, reflect.Float64:
		float, ok := yamlPrimitive.(float64)
		if !ok {

			integer, ok := yamlPrimitive.(int)
			if !ok {
				return fmt.Errorf("incompatible types go:(%s) yaml:(%s)", goPrimitive.Type().Name(), reflect.TypeOf(yamlPrimitive).Name())
			}

			goPrimitive.SetFloat(float64(integer))
			return nil
		}

		goPrimitive.SetFloat(float)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, ok := yamlPrimitive.(int)
		if !ok {
			return fmt.Errorf("incompatible types go:(%s) yaml:(%s)", goPrimitive.Type().Name(), reflect.TypeOf(yamlPrimitive).Name())
		}

		goPrimitive.SetInt(int64(integer))
		return nil
	case reflect.String:
		str, ok := yamlPrimitive.(string)
		if !ok {
			return fmt.Errorf("incompatible types go:(%s) yaml:(%s)", goPrimitive.Type().Name(), reflect.TypeOf(yamlPrimitive).Name())
		}

		goPrimitive.SetString(str)
		return nil
	case reflect.Bool:
		boolean, ok := yamlPrimitive.(bool)
		if !ok {
			return fmt.Errorf("incompatible types go:(%s) yaml:(%s)", goPrimitive.Type().Name(), reflect.TypeOf(yamlPrimitive).Name())
		}

		goPrimitive.SetBool(boolean)
		return nil
	case reflect.Interface:
		goPrimitive.Set(reflect.ValueOf(yamlPrimitive))
		return nil
	}

	return fmt.Errorf("unhandled types go:(%s) yaml:(%s)", goPrimitive.Type().Name(), reflect.TypeOf(yamlPrimitive).Name())
}

func fieldMap(structTag string, typ reflect.Type) map[string]reflect.StructField {
	fields := make(map[string]reflect.StructField)

	nFields := typ.NumField()
	for i := 0; i < nFields; i++ {
		field := typ.Field(i)

		name := field.Name
		if tagName, ok := field.Tag.Lookup(structTag); ok {
			parts := strings.Split(tagName, ",")
			if len(parts[0]) != 0 {
				name = parts[0]
			}
		}
		fields[name] = field
	}

	return fields
}
