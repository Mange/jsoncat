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

type fileresult struct {
	filename string
	index    int
	data     interface{}
	err      error
}

func ReadFiles(filenames []string) ([]interface{}, []error) {
	/*
		Read files concurrently while making sure that the data is added to the
		list in order.
		We also need to collect every error and return it separately.
	*/

	total := len(filenames)
	channel := make(chan fileresult, total)

	for index := range filenames {
		i := index
		go readFile(filenames[i], i, channel)
	}

	data := make([]interface{}, total)
	var errs []error
	totalRead := 0

	for result := range channel {
		if result.err == nil {
			data[result.index] = result.data
		} else {
			errs = append(errs, errors.New(fmt.Sprint(result.filename, ": ", result.err.Error())))
		}

		totalRead++

		if totalRead >= total {
			close(channel)
		}
	}

	return data, errs
}

func readFile(filename string, index int, channel chan fileresult) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		channel <- fileresult{filename, index, nil, err}
		return
	}

	doc, err := ReadJson(file)
	channel <- fileresult{filename, index, doc, err}
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

	data, errors := ReadFiles(flag.Args())

	if len(errors) > 0 {
		for _, err := range errors {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
		}
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
