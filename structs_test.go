package structs_test

import (
	"testing"

	"github.com/Nigel2392/go-structs"
)

func TestMakeStruct(t *testing.T) {
	var s = structs.New("json")
	s.StringField("Name", "name", true)
	s.IntField("Age", "age", true)
	s.BoolField("Is_cool", "is_cool", true)

	s.Make()

	s.SetField("Name", "Nigel")
	s.SetField("Age", 23)
	s.SetField("Is_cool", true)

	// Copy the struct
	var v = s.DeepCopy()

	if v.GetField("Name") != "Nigel" {
		t.Errorf("Expected %s, got %s", "Nigel", v.GetField("name"))
	}
	if v.GetField("Age") != 23 {
		t.Errorf("Expected %d, got %d", 23, v.GetField("age"))
	}
	if v.GetField("Is_cool") != true {
		t.Errorf("Expected %t, got %t", true, v.GetField("is_cool"))
	}

	v.SetField("Name", "Nigel2")
	v.SetField("Age", 24)
	v.SetField("Is_cool", false)
	// Change the new instance
	if v.GetField("Name") != "Nigel2" {
		t.Errorf("Expected %s, got %s", "Nigel2", v.GetField("name"))
	}
	if v.GetField("Age") != 24 {
		t.Errorf("Expected %d, got %d", 24, v.GetField("age"))
	}
	if v.GetField("Is_cool") != false {
		t.Errorf("Expected %t, got %t", false, v.GetField("is_cool"))
	}
	// Make sure the original instance is unchanged
	if s.GetField("Name") != "Nigel" {
		t.Errorf("Expected %s, got %s", "Nigel", s.GetField("name"))
	}
	if s.GetField("Age") != 23 {
		t.Errorf("Expected %d, got %d", 23, s.GetField("age"))
	}
	if s.GetField("Is_cool") != true {
		t.Errorf("Expected %t, got %t", true, s.GetField("is_cool"))
	}
}
