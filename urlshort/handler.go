package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.

// https://medium.com/geekculture/demystifying-http-handlers-in-golang-a363e4222756
// Read the above - it really helps in understanding the handler function
// There are three ways to serve a webpage using http.Handler:
// - Using the Handler naiively
// - Handlefunc: Allows you to pass a function
// - HandlerFunc: mixes and matches supporting to pass in a function and being able to add it to the handler
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	servePathToUrls := func(w http.ResponseWriter, r *http.Request) {
		//https://stackoverflow.com/questions/13533681/when-do-gos-pointers-dereference-themselves
		toRedirectUrl := pathsToUrls[r.URL.Path]

		fmt.Printf("After modification: %+v", r.URL.Path)
		fmt.Printf("Path Redirected: %+v", toRedirectUrl)
		http.Redirect(w, r, toRedirectUrl, http.StatusFound)
	}
	return http.HandlerFunc(servePathToUrls)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type urlPath struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYaml(data []byte) ([]urlPath, error) {
	var parsedData []urlPath

	err := yaml.Unmarshal(data, &parsedData)

	return parsedData, err
}

func buildMap(urlPaths []urlPath) map[string]string {
	m := make(map[string]string)
	for _, route := range urlPaths {
		m[route.Path] = route.Url
	}
	return m
}

// (http.HandlerFunc, error)
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yaml)
	fmt.Printf("%+v\n", parsedYaml)
	
	if err != nil {
		// fmt.Printf("%+v\n", err)
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	// fmt.Printf("%+v\n", pathMap)

	return MapHandler(pathMap, fallback), nil
}
