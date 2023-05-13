package structs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func IsRequired(field reflect.StructField) bool {
	return strings.Contains(field.Tag.Get("structs"), "required")
}

type Struct struct {
	// There is an optional parameter "required" for the fields of the struct.
	//
	// This can be used to determine whether the field is required or not in serialization for example.
	tag          string                // Default tag to use for enc_name
	fieldsByName []reflect.StructField // Inner fields.
	sstruct      reflect.Type          // The struct type
	structValue  reflect.Value         // The struct value
	made         bool                  // Whether the struct has been made or not
}

func From(v interface{}, tag string, fields ...string) *Struct {
	var s = New(tag)
	var structTyp reflect.Type
	var structVal reflect.Value
	switch v := v.(type) {
	case reflect.Value:
		if v.Kind() != reflect.Struct {
			panic(fmt.Sprintf("Cannot create struct from value of type %s", v.Kind().String()))
		}
		structVal = v
		structTyp = structVal.Type()
	case reflect.Type:
		if v.Kind() != reflect.Struct {
			panic(fmt.Sprintf("Cannot create struct from type of type %s", v.Kind().String()))
		}
		structTyp = v
		structVal = reflect.New(structTyp).Elem()
	case Struct:
		structTyp = v.sstruct
		structVal = reflect.New(structTyp).Elem()
	default:
		structTyp = reflect.TypeOf(v)
		structVal = reflect.New(structTyp).Elem()
	}
	if len(fields) > 0 {
		for _, field := range fields {
			var f, ok = structTyp.FieldByName(field)
			if !ok {
				panic(fmt.Sprintf("Field %s does not exist in struct %s", field, structTyp.Name()))
			}
			var enc_name = f.Tag.Get(tag)
			if enc_name == "-" {
				continue
			}
			s.AddField(field, enc_name, f.Type)
		}
		return s
	}
	for i := 0; i < structTyp.NumField(); i++ {
		var field = structTyp.Field(i)
		var enc_name = field.Tag.Get(tag)
		if enc_name == "-" {
			continue
		}
		var absolute_name = field.Name
		s.AddField(absolute_name, enc_name, field.Type)
	}
	return s
}

func New(tag string) *Struct {
	return &Struct{
		tag:          tag,
		fieldsByName: make([]reflect.StructField, 0),
	}
}

func (s *Struct) FieldByName(name string) reflect.Value {
	s.checkMade("Cannot get field by name if struct has not been made")
	return s.structValue.FieldByName(name)
}

// Field returns the field at the given index
//
// It will panic if the struct has not been made.
func (s *Struct) Field(index int) reflect.StructField {
	s.checkMade("Cannot get field by index if struct has not been made")
	return s.sstruct.Field(index)
}

func (s *Struct) checkMade(msg string) {
	if !s.made {
		panic(msg)
	}
}

// IsValid returns whether the struct is valid or not
//
// It will return false if the struct has not been made.
//
// Thus, it will return false if any fields have been added, and the made flag has been reset.
//
// Most other methods will panic if the struct is not made, or valid.
func (s *Struct) IsValid() bool {
	return s.made && s.structValue.IsValid()
}

// NumField returns the amount of fields in the struct
//
// If the struct has not been made, it will return 0.
func (s *Struct) NumField() int {
	if !s.made {
		return 0
	}
	return s.sstruct.NumField()
}

// Returns the amount of fields that have been added to the struct
//
// This is useful for when you want to know how many fields have been added, but the struct has not been made yet.
func (s *Struct) NumUninitializedField() int {
	return len(s.fieldsByName)
}

func (s *Struct) Interface() interface{} {
	s.checkMade("Cannot get interface if struct has not been made")
	return s.structValue.Interface()
}

func (s *Struct) PtrTo() interface{} {
	s.checkMade("Cannot get pointer to if struct has not been made")
	return s.structValue.Addr().Interface()
}

func (s *Struct) SetField(name string, value interface{}) {
	s.checkMade("Cannot set field if struct has not been made")
	var field = s.structValue.FieldByName(name)
	if !field.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist", name))
	}
	var valueOf = valueOf(value)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if field.Kind() != valueOf.Kind() {
		panic(fmt.Sprintf("Cannot set field %s with value of type %s", name, valueOf.Kind().String()))
	}
	field.Set(valueOf)
}

func (s *Struct) SetFieldByIndex(index int, value interface{}) {
	s.checkMade("Cannot set field if struct has not been made")
	var field = s.structValue.Field(index)
	if !field.IsValid() {
		panic(fmt.Sprintf("Field %d does not exist", index))
	}
	var valueOf = valueOf(value)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if field.Kind() != valueOf.Kind() {
		panic(fmt.Sprintf("Cannot set field %d with value of type %s", index, valueOf.Kind().String()))
	}
	field.Set(valueOf)
}

// Deep copy of the struct
//
// This is useful for when you want to modify a struct without modifying the original
func (s *Struct) DeepCopy() *Struct {
	s.checkMade("Cannot deep copy if struct has not been made")
	var newStruct = New(s.tag)
	for _, field := range s.fieldsByName {
		newStruct.AddField(field.Name, field.Tag.Get(s.tag), field.Type, field.Tag.Get("structs") == "required")
	}

	newStruct.Make()

	for _, field := range s.fieldsByName {
		var newFieldByIndex = newStruct.structValue.FieldByIndex(field.Index)
		if newFieldByIndex.Kind() == reflect.Ptr {
			newFieldByIndex = newFieldByIndex.Elem()
		}
		var fieldByIndex = s.structValue.FieldByIndex(field.Index)
		if fieldByIndex.Kind() == reflect.Ptr {
			fieldByIndex = fieldByIndex.Elem()
		}

		if fieldByIndex.Kind() != newFieldByIndex.Kind() {
			panic(fmt.Sprintf("Cannot deep copy field %s, because the types are different", field.Name))
		}

		if !newFieldByIndex.CanSet() {
			panic(fmt.Sprintf("Cannot deep copy field %s, because it cannot be set", field.Name))
		}

		newFieldByIndex.Set(fieldByIndex)
	}
	return newStruct
}

func valueOf(v interface{}) reflect.Value {
	switch v.(type) {
	case reflect.Value:
		return v.(reflect.Value)
	default:
		return reflect.ValueOf(v)
	}
}

func (s *Struct) AddField(absolute_name, enc_name string, typeOf reflect.Type, required ...bool) {
	if absolute_name == "" {
		panic("Field name cannot be empty")
	}
	if enc_name == "" {
		enc_name = absolute_name
	}
	for _, field := range s.fieldsByName {
		if field.Name == absolute_name {
			panic(fmt.Sprintf("Field %s already exists", absolute_name))
		}
	}

	// If the struct has already been made,
	// we need to reset the flag so the Make() method will re-make it
	s.made = false
	var tag string = fmt.Sprintf(`%s:"%s"`, s.tag, enc_name)
	if len(required) > 0 && required[0] {
		tag += fmt.Sprintf(` structs:"required"`)
	}
	var field = reflect.StructField{
		Name:      absolute_name,
		Tag:       reflect.StructTag(tag),
		Type:      typeOf,
		Anonymous: false,
	}
	s.fieldsByName = append(s.fieldsByName, field)
}

func (s *Struct) AddStructField(field reflect.StructField) {
	if field.Name == "" {
		panic("Field name cannot be empty")
	}
	for _, f := range s.fieldsByName {
		if f.Name == field.Name {
			panic(fmt.Sprintf("Field %s already exists", field.Name))
		}
	}
	if field.Anonymous {
		panic("Cannot add anonymous field")
	}

	// If the struct has already been made,
	// we need to reset the flag so the Make() method will re-make it
	s.made = false
	s.fieldsByName = append(s.fieldsByName, field)
}

func (s *Struct) StringField(absolute_name, name string, required ...bool) {
	s.AddField(absolute_name, name, reflect.TypeOf(""), required...)
}

func (s *Struct) IntField(absolute_name, name string, required ...bool) {
	s.AddField(absolute_name, name, reflect.TypeOf(0), required...)
}

func (s *Struct) FloatField(absolute_name, name string, required ...bool) {
	s.AddField(absolute_name, name, reflect.TypeOf(0.0), required...)
}

func (s *Struct) BoolField(absolute_name, name string, required ...bool) {
	s.AddField(absolute_name, name, reflect.TypeOf(false), required...)
}

func (s *Struct) SliceField(absolute_name, name string, typeOf reflect.Type, required ...bool) {
	s.AddField(absolute_name, name, reflect.SliceOf(typeOf), required...)
}

func (s *Struct) MapField(absolute_name, name string, typeOfKey, typeOfValue reflect.Type, required ...bool) {
	if !typeOfKey.Comparable() {
		panic(fmt.Sprintf("Map key type %s is not comparable", typeOfKey.String()))
	}
	s.AddField(absolute_name, name, reflect.MapOf(typeOfKey, typeOfValue), required...)
}

func (s *Struct) StructField(absolute_name, name string, other *Struct, required ...bool) {
	s.AddField(absolute_name, name, other.sstruct, required...)
}

func (s *Struct) GetField(name string) interface{} {
	s.checkMade("Cannot get field if struct has not been made")
	var field = s.structValue.FieldByName(name)
	if !field.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist", name))
	}
	return field.Interface()
}

func (s *Struct) Remake() {
	s.made = false
	s.Make()
}

func (s *Struct) Make() {
	if !s.made {
		s.sstruct = reflect.StructOf(s.fieldsByName)
		s.made = true
	}
	if s.made {
		var NewOf = reflect.New(s.sstruct)
		s.structValue = NewOf.Elem()
	}
}

func ScanInto(s, dest interface{}, fields ...string) error {
	var iFace any
	switch s.(type) {
	case *Struct:
		iFace = s.(*Struct).Interface()
	default:
		iFace = s
	}
	return scanInto(iFace, dest, fields...)
}

func scanInto(s, dest interface{}, fields ...string) error {
	var typeOfSource = reflect.TypeOf(s)
	var valueOfSource = reflect.ValueOf(s)
	var typeOfDest = reflect.TypeOf(dest)
	if typeOfSource.Kind() == reflect.Ptr {
		typeOfSource = typeOfSource.Elem()
		valueOfSource = valueOfSource.Elem()
	}
	if typeOfSource.Kind() != reflect.Struct {
		return fmt.Errorf("Source is not a struct")
	}
	if typeOfDest.Kind() != reflect.Ptr {
		return fmt.Errorf("Destination is not a pointer")
	}
	if typeOfDest.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Destination is not a pointer to a struct")
	}
	var valueOfDest = reflect.ValueOf(dest)
	var valueOfDestElem = valueOfDest.Elem()
	var numFields = typeOfSource.NumField()
	for i := 0; i < numFields; i++ {
		var field = typeOfSource.Field(i)
		var name = field.Name
		if len(fields) > 0 {
			var found bool
			for _, f := range fields {
				if f == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		var value = valueOfSource.Field(i)
		var destField = valueOfDestElem.FieldByName(name)
		if !destField.IsValid() {
			continue
		}
		if !destField.CanSet() {
			continue
		}
		if destField.Type() != value.Type() {
			continue
		}
		destField.Set(value)
	}
	return nil
}

func (s *Struct) MarshalJSON() ([]byte, error) {
	s.checkMade("Cannot marshal if struct has not been made")
	return json.Marshal(s.structValue.Interface())
}

func (s *Struct) UnmarshalJSON(data []byte) error {
	s.checkMade("Cannot unmarshal if struct has not been made")
	return json.Unmarshal(data, s.structValue.Addr().Interface())
}
