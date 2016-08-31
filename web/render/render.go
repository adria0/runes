package render

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/adriamb/gopad/store"
	"github.com/russross/blackfriday"
)

const (
	commonHTMLFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	commonExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

type renderer interface {
	BlockDescriptor() string
}

type divRenderer interface {
	renderer
	HTMLHeaders() string
	RenderToBuffer(data string, params string) (string, error)
}

type imageRenderer interface {
	renderer
	ImageFileExtension() string
	RenderToFile(data string, params string, filename string) error
}

var (
	renderers         map[string]renderer
	errNotImplemented = errors.New("Not implemented")
)

func init() {

	rs := [...]renderer{
		&dotRenderer{},
		&goatRenderer{},
		&jsseqRenderer{},
		&jsflowRenderer{},
	}

	renderers = make(map[string]renderer)
	for _, r := range rs {
		renderers[r.BlockDescriptor()] = r
	}

}

func params2map(params string, defaults map[string]string) map[string]string {

	ret := make(map[string]string)
	for k, v := range defaults {
		ret[k] = v
	}

	split := strings.Split(params, ";")
	for _, param := range split {
		if strings.Contains(param, "=") {
			pair := strings.Split(param, "=")
			ret[pair[0]] = pair[1]
		} else {
			ret[param] = "true"
		}
	}

	return ret
}

func mustWrite(w io.Writer, p []byte) {
	_, err := w.Write(p)
	if err != nil {
		log.Panic(err)
	}
}

func mustWriteString(w io.Writer, s string) {
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Panic(err)
	}
}

func blockRenderer(content string, srange blackfriday.SourceRange, langAndParams string) (string, error) {

	lap := strings.Split(langAndParams, "|")
	language := lap[0]

	var class = ""
	var params = ""

	if len(lap) > 1 {
		class = lap[1]
	}
	if len(lap) > 2 {
		params = lap[2]
	}

	var r renderer
	var exists bool

	if r, exists = renderers[language]; !exists {
		return "", fmt.Errorf("Handle for language " + language + " does not exist")
	}

	if imgR, ok := r.(imageRenderer); ok {

		hasher := sha1.New()
		mustWriteString(hasher, content)
		sha1 := hasher.Sum(nil)
		ID := hex.EncodeToString(sha1[:])
		filename := ID + "." + imgR.ImageFileExtension()

		if !store.ExistsCache(filename) {

			filenameWithPath := store.GetCachePath(filename)

			if err := imgR.RenderToFile(content, params, filenameWithPath); err != nil {
				return "", err
			}

		}

		imgloc := fmt.Sprintf("<img src=/cache/%s class=\""+class+"\" %s><br>", filename, srange.Attrs())

		return imgloc, nil

	}

	divR := r.(divRenderer)

	rendered, err := divR.RenderToBuffer(content, params)
	if err != nil {
		return "", err
	}

	div := fmt.Sprintf("<div class=\""+class+"\" %s>%s</div><br>", srange.Attrs(), string(rendered))

	return div, nil

}

// HTMLHeaders requiered in the html page to run the plugins
func HTMLHeaders() string {

	var buffer bytes.Buffer

	for _, r := range renderers {

		if divR, ok := r.(divRenderer); ok {

			_, err := buffer.WriteString(divR.HTMLHeaders())
			if err != nil {
				panic(err)
			}
		}

	}

	return buffer.String()
}

// Render a markdown into html
func Render(markdown string) []byte {

	params := blackfriday.HtmlRendererParameters{
		BlockRenderer: blockRenderer,
	}

	renderer := blackfriday.HtmlRendererWithParameters(commonHTMLFlags, "", "", params)

	html := blackfriday.Markdown([]byte(markdown), renderer, commonExtensions)

	var out bytes.Buffer
	mustWriteString(&out, "<div class='markdown'>")
	mustWrite(&out, []byte(html))
	mustWriteString(&out, "</div>")
	return out.Bytes()
}
