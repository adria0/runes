package render

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/adriamb/gopad/dict"
	"github.com/adriamb/gopad/store"
	"github.com/russross/blackfriday"
)

type renderHandler struct {
	filename func(string) string
	render   func(string, string, []byte) error
}

const (
	extPNG = ".png"
)

var (
	renderHandlers = map[string]renderHandler{
		"dot":   {filenameDot, renderDot},
		"umlet": {filenameUmlet, renderUmlet},
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

// Render a markdown into html
func Render(markdown string, dict *dict.Dict) []byte {

	rendered := renderImages(markdown)
	html := string(blackfriday.MarkdownCommon(rendered))

	if defs, err := dict.Defs(); err == nil {
		for k, v := range defs {
			v = `<a href="#"><span title="` + v + `">` + k + `</span></a>`
			html = strings.Replace(html, "§"+k, v, -1)
		}
	} else {
		log.Print("Render error", err)
		return []byte("render error")
	}

	var out bytes.Buffer
	mustWriteString(&out, "<div class='markdown'>")
	mustWrite(&out, []byte(html))
	mustWriteString(&out, "</div>")
	return out.Bytes()
}

func renderImages(markdown string) []byte {
	var out bytes.Buffer
	var block bytes.Buffer
	var handler renderHandler
	var imagetags string
	var params string

	lines := strings.Split(markdown, "\n")
	inblock := false
	writer := &out

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inblock = !inblock
			if inblock {
				blocktype := strings.TrimSpace(line[3:])
				if strings.Contains(blocktype, ";") {
					split := strings.Split(blocktype, ";")
					blocktype = split[0]
					imagetags = split[1]
				} else {
					imagetags = ""
				}
				if strings.Contains(blocktype, ":") {
					split := strings.Split(blocktype, ":")
					blocktype = split[0]
					params = split[1]
				} else {
					params = ""
				}
				var exists bool
				if handler, exists = renderHandlers[blocktype]; exists {
					block = bytes.Buffer{}
					writer = &block
					continue
				}
			} else {
				if writer == &block {
					writer = &out
					data := block.Bytes()
					hasher := sha1.New()
					mustWrite(hasher, data)
					mustWriteString(hasher, params)
					sha1 := hasher.Sum(nil)
					ID := hex.EncodeToString(sha1[:])
					filename := handler.filename(ID)
					if !store.ExistsCache(filename) {
						if err := handler.render(filename, params, data); err != nil {
							mustWriteString(&out, fmt.Sprintf("%v", err))
							continue
						}
					}
					mustWriteString(&out, "!["+imagetags+"](/cache/"+filename+")\n")
					continue
				}
			}
		}
		mustWriteString(writer, line)
		mustWriteString(writer, "\n")
	}

	if inblock {
		return []byte("Unaligned blocks")
	}

	return out.Bytes()
}
