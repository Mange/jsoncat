package main

type docType int

const (
	arrays docType = iota
	objects
	bare
	mixed
)

func MergeJson(docs []interface{}) interface{} {
	types := determineTypes(docs)
	switch types {
	case arrays:
		return mergeJsonArrays(docs)
	case objects:
		return mergeJsonObjects(docs)
	default:
		return nil
	}
}

func determineTypes(docs []interface{}) docType {
	hasArrays := false
	hasObjects := false
	hasBare := false

	for _, doc := range docs {
		switch doc.(type) {
		case []interface{}:
			hasArrays = true
		case map[string]interface{}:
			hasObjects = true
		default:
			hasBare = true
		}
	}

	if hasArrays && hasObjects || hasArrays && hasBare || hasObjects && hasBare {
		return mixed
	} else if hasArrays {
		return arrays
	} else if hasObjects {
		return objects
	} else {
		return bare
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
