package mergo_test

import (
	"reflect"
	"testing"

	"github.com/imdatngo/mergo"
)

// Embedded to test embedded struct
type Embedded struct {
	InnerField int `json:"inner_field"`
	FindMe     string
	FieldName  string `json:"field_name"`
}

type jsonTest struct {
	*Embedded
	FieldName    string `json:"field_name"`
	JSONWithType string `json:"another_field,string"`
	NoJSONTag    string
	EmptyJSONTag string  `json:",string"`
	DoNotMerge   string  `json:"-"`
	Pointer      *string `json:"pointer,omitempty"`
}

func TestMap2Struct_WithJSONTagLookup(t *testing.T) {
	teststr := "string ptr"
	scrMap := map[string]interface{}{
		"field_name":    "FieldName",
		"another_field": "JSONWithType",
		"NoJSONTag":     "no_json_tag",
		"EmptyJSONTag":  "empty",
		"DoNotMerge":    "this should not be merged",
		"-":             "tricky",
		"non":           "exist",
		"pointer":       &teststr,
		"inner_field":   123,
		"FindMe":        "yay!",
	}
	dstStruct := jsonTest{}
	expectStruct := jsonTest{
		Embedded: &Embedded{
			InnerField: 123,
			FindMe:     "yay!",
		},
		FieldName:    "FieldName",
		JSONWithType: "JSONWithType",
		NoJSONTag:    "no_json_tag",
		EmptyJSONTag: "empty",
		DoNotMerge:   "",
		Pointer:      &teststr,
	}

	if err := mergo.Map(&dstStruct, scrMap, mergo.WithJSONTagLookup); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(dstStruct, expectStruct) {
		t.Fatalf("map -> struct failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", dstStruct, expectStruct)
	}
}

func TestStruct2Map_WithJSONTagLookup(t *testing.T) {
	srcStruct := jsonTest{
		Embedded: &Embedded{
			InnerField: 123,
			FindMe:     "yay!",
			FieldName:  "this will be ignore because of conflict with outer struct",
		},
		FieldName:    "FieldName",
		JSONWithType: "JSONWithType",
		NoJSONTag:    "no_json_tag",
		EmptyJSONTag: "empty",
		DoNotMerge:   "this should not be merged",
	}
	dstMap := make(map[string]interface{})
	expectMap := map[string]interface{}{
		"inner_field":   123,
		"FindMe":        "yay!",
		"field_name":    "FieldName",
		"another_field": "JSONWithType",
		"NoJSONTag":     "no_json_tag",
		"EmptyJSONTag":  "empty",
	}

	if err := mergo.Map(&dstMap, srcStruct, mergo.WithJSONTagLookup); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(dstMap, expectMap) {
		t.Fatalf("struct -> map failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", dstMap, expectMap)
	}
}
