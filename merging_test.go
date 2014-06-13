package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func parseJson(data string) interface{} {
	var unmarshaled interface{}
	err := json.Unmarshal([]byte(data), &unmarshaled)
	check(err)
	return unmarshaled
}

func expectEqual(t *testing.T, methodName string, in, expected, actual interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s(%+v) = %+v, but should be %+v", methodName, in, actual, expected)
	}
}

func expectMergeJson(t *testing.T, in []interface{}, expected interface{}) {
	actual, err := MergeJson(in)
	if err != nil {
		t.Error(err)
	}
	expectEqual(t, "MergeJson", in, expected, actual)
}

func expectMergeJsonToFail(t *testing.T, in []interface{}, message string) {
	actual, err := MergeJson(in)

	if err == nil {
		t.Errorf("No error was returned! Returned data: %v", actual)
	}

	if err.Error() != message {
		t.Errorf("Expected error message to be \"%s\", was: \"%s\"", message, err.Error())
	}
}

func TestMergeJson_arrays(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("[1, 2, 3]")
	in[1] = parseJson("[4, 5, 6]")

	expected := parseJson("[1, 2, 3, 4, 5, 6]")

	expectMergeJson(t, in, expected)
}

func TestMergeJson_arrays_of_different_types(t *testing.T) {
	in := make([]interface{}, 3)
	in[0] = parseJson("[1337]")
	in[1] = parseJson("[\"hello\"]")
	in[2] = parseJson("[\"world\"]")

	expected := parseJson("[1337, \"hello\", \"world\"]")

	expectMergeJson(t, in, expected)
}

func TestMergeJson_objects(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("{\"a\": 1, \"b\": 5}")
	in[1] = parseJson("{\"a\": 2}")

	expected := parseJson("{\"a\": 2, \"b\": 5}")

	expectMergeJson(t, in, expected)
}

func TestMergeJson_objects_will_not_deep_merge(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("{\"a\": [1, 2, 3]}")
	in[1] = parseJson("{\"a\": [3, 4, 5]}")

	expected := parseJson("{\"a\": [3, 4, 5]}")

	expectMergeJson(t, in, expected)
}

func TestMergeJson_with_mixed_types(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("{\"a\": 1}")
	in[1] = parseJson("[3, 4, 5]")

	expectMergeJsonToFail(t, in, "Cannot merge incompatible types: Object, Array")
}
