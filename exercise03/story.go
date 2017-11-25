package cyoa

import (
	"encoding/json"
	"io"
)

// JSONStory parses a JSON story file into a Story object.
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story consists of a map of dynamic story names to a Chapter struct.
type Story map[string]Chapter

// Chapter represents a portion of a Story that has a Title, Paragraphs, and
// (optional) Options for what to do at the end of the Chapter.
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option represents a link to the next Chapter in the Story, if one exist.
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
