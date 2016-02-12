package render

import (
    "bytes"
    "strings"
    "fmt"
    "crypto/sha1"
    "os/exec"
    "encoding/hex"
	"github.com/russross/blackfriday"
    "github.com/amassanet/gopad/store"
)

type renderHandler struct {
    filename func(string) string
    render func(string,[]byte) error
}

const (
    extPNG = ".png"
)

var (
    renderHandlers = map[string]renderHandler{
        "dot" : {filenameDot,renderDot },
    }
)

func Render(markdown string) []byte {
    rendered :=  renderImages(markdown)
    return blackfriday.MarkdownCommon(rendered)
}


func renderImages(markdown string) []byte {
    var out  bytes.Buffer
    var block  bytes.Buffer
    var handler renderHandler
    var imagetags string

    lines := strings.Split(markdown,"\n")
    inblock := false
    writer := &out

    for _,line := range lines {
        if strings.HasPrefix(line,"```") {
            inblock = !inblock
            if inblock {
                blocktype := strings.TrimSpace(line[3:])
                if strings.Contains(blocktype,";") {
                    split := strings.Split(blocktype,";")
                    blocktype = split[0]
                    imagetags = split[1]
                } else {
                    imagetags = ""
                }
                var exists bool
                if handler, exists = renderHandlers[blocktype] ; exists {
                    block = bytes.Buffer{}
                    writer = &block
                    continue
                }
            } else {
                if writer == &block {
                    writer = &out
                    data := block.Bytes()
                    sha1 := sha1.Sum(data)
                    ID := hex.EncodeToString(sha1[:])
                    filename := handler.filename(ID)
                    if !store.ExistsCache(filename) {
                        if err := handler.render(filename,data); err!=nil {
                            out.WriteString(fmt.Sprintf("%v",err))
                            continue
                        }
                    }
                    out.WriteString("!["+imagetags+"](/cache/"+filename+")\n")
                    continue
                }
            }
        }
        writer.WriteString(line)
        writer.WriteString("\n")
    }

    if inblock {
       return []byte("Unaligned blocks")
    }

    return out.Bytes()
}

func filenameDot(ID string) string {
    return ID+".png"
}

func renderDot(filename string, data []byte) error {
    var dot  bytes.Buffer
    dot.WriteString("digraph G {\n")
    dot.Write(data)
    dot.WriteString("\n}")

    dotfile,err := store.WriteTemp(filename+".dot",dot.Bytes())
    if err != nil {
        return fmt.Errorf("%v",err)
    }
    pngfile := store.GetCachePath(filename)

    cmd := exec.Command("dot",dotfile,"-Tpng","-o"+pngfile)
    var out bytes.Buffer
    cmd.Stderr = &out

    if err := cmd.Run() ; err != nil {
        return fmt.Errorf("%v %v",err,out.String())
    }
    return nil
}

