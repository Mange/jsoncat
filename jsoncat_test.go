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

func TestMergeJson_arrays(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("[1, 2, 3]")
	in[1] = parseJson("[4, 5, 6]")

	expected := parseJson("[1, 2, 3, 4, 5, 6]")

	expectEqual(t, "MergeJson", in, expected, MergeJson(in))
}

func TestMergeJson_arrays_of_different_types(t *testing.T) {
	in := make([]interface{}, 3)
	in[0] = parseJson("[1337]")
	in[1] = parseJson("[\"hello\"]")
	in[2] = parseJson("[\"world\"]")

	expected := parseJson("[1337, \"hello\", \"world\"]")

	expectEqual(t, "MergeJson", in, expected, MergeJson(in))
}

func TestMergeJson_objects(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("{\"a\": 1, \"b\": 5}")
	in[1] = parseJson("{\"a\": 2}")

	expected := parseJson("{\"a\": 2, \"b\": 5}")

	expectEqual(t, "MergeJson", in, expected, MergeJson(in))
}

func TestMergeJson_objects_will_not_deep_merge(t *testing.T) {
	in := make([]interface{}, 2)
	in[0] = parseJson("{\"a\": [1, 2, 3]}")
	in[1] = parseJson("{\"a\": [3, 4, 5]}")

	expected := parseJson("{\"a\": [3, 4, 5]}")

	expectEqual(t, "MergeJson", in, expected, MergeJson(in))
}
