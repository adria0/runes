package render

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/adriamb/gopad/store"
	"github.com/russross/blackfriday"
)

type renderHandler struct {
	filename func(string) string
	render   func(string, string, []byte) error
}

const (
	extPNG = ".png"

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

var (
	renderHandlers = map[string]renderHandler{
		"dot":      {filenameDot, renderDot},
	}
)

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

func blockRenderer(content []byte, srange blackfriday.SourceRange, langAndParams string) ([]byte, error) {

	var handler renderHandler
	var exists bool

	lap := strings.Split(langAndParams, ":")
	language := lap[0]

	var class = ""
	if len(lap) > 1 {
		class = lap[1]
	}

	if handler, exists = renderHandlers[language]; !exists {
		return nil, fmt.Errorf("Handle for language " + language + " does not exist")
	}

	hasher := sha1.New()
	mustWrite(hasher, content)
	sha1 := hasher.Sum(nil)
	ID := hex.EncodeToString(sha1[:])
	filename := handler.filename(ID)

	if !store.ExistsCache(filename) {
		if err := handler.render(filename, "", content); err != nil {
			return nil, err
		}
	}

	imgloc := fmt.Sprintf("<img src=/cache/%s class=\""+class+"\" %s><br>", filename, srange.Attrs())
	imglocbytes := []byte(imgloc)

	return imglocbytes, nil
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
