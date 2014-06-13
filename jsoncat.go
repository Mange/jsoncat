package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ReadFiles(filenames []string) []interface{} {
	data := make([]interface{}, len(filenames))

	for index, filename := range filenames {
		file, err := os.Open(filename)
		check(err)
		data[index] = ReadJson(file)
	}

	return data
}

func ReadJson(r io.Reader) interface{} {
	var doc interface{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&doc)
	check(err)
	return doc
}

func WriteJsonDocuments(docs []interface{}, w io.Writer) {
	var err error
	encoder := json.NewEncoder(w)
	if len(docs) == 1 {
		err = encoder.Encode(docs[0])
	} else {
		err = encoder.Encode(docs)
	}
	check(err)
}

func WriteJsonDocument(doc interface{}, w io.Writer) {
	var err error
	encoder := json.NewEncoder(w)
	err = encoder.Encode(doc)
	check(err)
}

func MergeJson(docs []interface{}) interface{} {
	return MergeJsonArrays(docs)
}

func MergeJsonArrays(docs []interface{}) []interface{} {
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

func main() {
	var merge bool

	flag.BoolVar(&merge, "merge", false, "Merge files")
	flag.Parse()

	data := ReadFiles(flag.Args())
	if merge {
		WriteJsonDocument(MergeJson(data), os.Stdout)
	} else {
		WriteJsonDocuments(data, os.Stdout)
	}
}
