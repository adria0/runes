package render

import (
	"bytes"
	"fmt"
	"github.com/adriamb/gopad/store"
	"os/exec"
)

func filenameDot(ID string) string {
	return ID + ".png"
}

func renderDot(filename string, params string, data []byte) error {
	var dot bytes.Buffer
	dot.WriteString("digraph G {\n")
	dot.Write(data)
	dot.WriteString("\n}")

	dotfile, err := store.WriteTemp(filename+".dot", dot.Bytes())
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	pngfile := store.GetCachePath(filename)

	cmd := exec.Command("dot", dotfile, "-Tpng", "-o"+pngfile)
	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v %v", err, out.String())
	}
	return nil
}
