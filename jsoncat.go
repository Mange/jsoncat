package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ReadFiles(filenames []string) ([]interface{}, error) {
	data := make([]interface{}, len(filenames))

	for index, filename := range filenames {
		file, err := os.Open(filename)
		check(err)
		doc, err := ReadJson(file)
		if err != nil {
			return nil, errors.New(fmt.Sprint(filename, ": ", err.Error()))
		}
		data[index] = doc
	}

	return data, nil
}

func ReadJson(r io.Reader) (interface{}, error) {
	var doc interface{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&doc)
	return doc, err
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

	data, err := ReadFiles(flag.Args())

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

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
