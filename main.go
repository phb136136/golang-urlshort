package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"golang-urlshort/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	yamlFile, err := ioutil.ReadFile("file.yaml")
		if err != nil {
		panic(err)
		}
	
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// fmt.Printf("+%v", mapHandler)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	// 	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
	// http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
