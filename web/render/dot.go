package render

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/adriamb/gopad/store"
)

func filenameDot(ID string) string {
	return ID + ".png"
}

func renderDot(filename string, params string, data []byte) error {
	var dot bytes.Buffer
	var err error

	if _, err = dot.WriteString("digraph G {\n"); err != nil {
		return err
	}
	if _, err = dot.Write(data); err != nil {
		return err
	}
	if _, err = dot.WriteString("\n}"); err != nil {
		return err
	}

	dotfile, err := store.WriteTemp(filename+".dot", dot.Bytes())
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	svgfile := store.GetCachePath(filename)

	cmd := exec.Command("dot", dotfile, "-Tpng", "-o"+svgfile)
	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v %v", err, out.String())
	}
	return nil
}
