package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToURLs := buildMap(pathURLs)
	return MapHandler(pathsToURLs, fallback), nil
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func parseYaml(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, path := range pathURLs {
		pathsToUrls[path.Path] = path.URL
	}
	return pathsToUrls
}

// JSONHandler will parse the provided JSON and then return an http.HandlerFunc
// (which also implements http.Handler) that will attempt to map any paths to
// their corresponding URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}
	pathsToURLs := buildMap(pathURLs)
	return MapHandler(pathsToURLs, fallback), nil
}

func parseJSON(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := json.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}
