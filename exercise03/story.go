package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var tpl *template.Template

var defaultHandlerTemplate = `
<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<title>Choose Your Own Adventure</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}
			<ul>
			{{range .Options}}
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
			</ul>
		</body>
	</html>`

// NewHandler returns an instance of a handler struct that contains a Story for
// the http.Handler to serve.
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}

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
