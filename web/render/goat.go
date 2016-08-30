package render

import (
	"strings"
	"os"

	"github.com/blampe/goat/src"
)

type goatRenderer struct {

}

func (g *goatRenderer) BlockDescriptor() string {
    return "goat"
}

func (g *goatRenderer) ImageFileExtension() string {
	return ".svg"
}

func (g *goatRenderer) RenderToFile(data string, params string, filename string) error{

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	goat.ASCIItoSVG(strings.NewReader(data), file)

	return nil
}
