package main

import (
	"fmt"
	"strings"
)

type docType int

const (
	array docType = iota
	object
	bare
	mixed
)

func (t docType) Name() string {
	switch t {
	case array:
		return "Array"
	case object:
		return "Object"
	case bare:
		return "Bare"
	case mixed:
		return "Mixed"
	default:
		return "Unknown type"
	}
}

type MergeError struct {
	message string
}

func MergeJson(docs []interface{}) (interface{}, error) {
	commonType, allTypes := determineTypes(docs)

	switch commonType {
	case array:
		return mergeJsonArrays(docs), nil
	case object:
		return mergeJsonObjects(docs), nil
	case bare:
		// TODO
	case mixed:
		return nil, newMixedError(allTypes)
	}

	return nil, newMergeError("Unexpected error")
}

func determineTypes(docs []interface{}) (docType, []docType) {
	hasArrays := false
	hasObjects := false
	hasBare := false
	allTypes := make([]docType, len(docs))

	for i, doc := range docs {
		switch doc.(type) {
		case []interface{}:
			hasArrays = true
			allTypes[i] = array
		case map[string]interface{}:
			hasObjects = true
			allTypes[i] = object
		default:
			hasBare = true
			allTypes[i] = bare
		}
	}

	if hasArrays && hasObjects || hasArrays && hasBare || hasObjects && hasBare {
		return mixed, allTypes
	} else if hasArrays {
		return array, allTypes
	} else if hasObjects {
		return object, allTypes
	} else {
		return bare, allTypes
	}
}

func mergeJsonArrays(docs []interface{}) []interface{} {
	totalElements := 0

	for _, doc := range docs {
		totalElements += len(doc.([]interface{}))
	}

	array := make([]interface{}, totalElements)

	combinedIndex := 0
	for _, doc := range docs {
		for _, value := range doc.([]interface{}) {
			array[combinedIndex] = value
			combinedIndex += 1
		}
	}

	return array
}

func mergeJsonObjects(docs []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, doc := range docs {
		object := doc.(map[string]interface{})
		for key, value := range object {
			result[key] = value
		}
	}
	return result
}

func newMergeError(message string) MergeError {
	return MergeError{message}
}

func newMixedError(types []docType) MergeError {
	typeNames := make([]string, len(types))
	for i, t := range types {
		typeNames[i] = t.Name()
	}
	message := fmt.Sprint("Cannot merge incompatible types: ", strings.Join(typeNames, ", "))
	return newMergeError(message)
}

func (e MergeError) Error() string {
	return e.message
}
