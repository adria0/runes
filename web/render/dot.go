package render

import (
	"bytes"
	"fmt"
	"github.com/adriamb/gopad/store"
	"os/exec"
)

type dotRenderer struct {
}

func (d *dotRenderer) BlockDescriptor() string {
	return "dot"
}

func (d *dotRenderer) ImageFileExtension() string {
	return "png"
}

func (d *dotRenderer) RenderToFile(data string, params string, filename string) error {

	var dot bytes.Buffer
	var err error

	if _, err = dot.WriteString("digraph G {\n"); err != nil {
		return err
	}
	if _, err = dot.WriteString(data); err != nil {
		return err
	}
	if _, err = dot.WriteString("\n}"); err != nil {
		return err
	}

	dotfile, err := store.WriteTemp(dot.Bytes())
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	cmd := exec.Command("dot", dotfile, "-Tpng", "-o"+filename)
	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v %v", err, out.String())
	}
	return nil

}
