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

func main() {
	var merge bool

	flag.BoolVar(&merge, "merge", false, "Merge files")
	flag.Parse()

	data := ReadFiles(flag.Args())
	if merge {
		data, err := MergeJson(data)
		if err == nil {
			WriteJsonDocument(data, os.Stdout)
		} else {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
			os.Exit(1)
		}
	} else {
		WriteJsonDocuments(data, os.Stdout)
	}
}
