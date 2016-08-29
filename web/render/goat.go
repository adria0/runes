package render

import (
	"bytes"
    "os"

    "github.com/blampe/goat/src"
	"github.com/adriamb/gopad/store"
)
func filenameGoat(ID string) string {
	return ID + ".svg"
}

func renderGoat(filename string, params string, data []byte) error {

	svgfile := store.GetCachePath(filename)
    file, err := os.Create(svgfile)
    if err != nil {
        return err
    }
    defer file.Close()
    goat.ASCIItoSVG(bytes.NewReader(data), file)

	return nil
}
