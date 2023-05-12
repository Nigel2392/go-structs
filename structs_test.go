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
}
