package render

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"

	"github.com/adriamb/gopad/store"
)

func filenameUmlet(ID string) string {
	return ID + ".png"
}

type xmlTemplate struct {
	Type string
	X    int
	Y    int
	W    int
	H    int
	Code string
}

var (
	diagramTypes = map[string]string{
		"sequence": "UMLSequenceAllInOne",
		"class":    "UMLClass",
	}
)

const xmltemplate = `
<diagram program="umlet" version="14.1.1">
<zoom_level>10</zoom_level>
<element>
<id>{{ .Type }}</id>
<coordinates>
<x>{{ .X }}</x>
<y>{{ .Y }}</y>
<w>{{ .W }}</w>
<h>{{ .H }}</h>
</coordinates>
<panel_attributes>{{ .Code  }}</panel_attributes>
<additional_attributes/>
</element>
</diagram>
`

func renderUmlet(filename string, params string, data []byte) error {

	umltype, found := diagramTypes[params]
	if !found {
		return fmt.Errorf("Bad type " + params)
	}

	tparams := xmlTemplate{
		Type: umltype,
		X:    0, Y: 0, W: 200, H: 400,
		Code: string(data),
	}

	tmpl, err := template.New("test").Parse(xmltemplate)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, tparams)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	uxffile, err := store.WriteTemp("temp.uxf", b.Bytes())
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	pngfile := store.GetCachePath(filename)

	cmd := exec.Command(
		"java", "-jar", "external/umlet/umlet.jar",
		"-action=convert", "-format=png",
		"-filename="+uxffile, "-output="+pngfile,
	)

	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v %v", err, out.String())
	}
	return nil

}
