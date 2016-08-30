package render

import (
	"fmt"
	"strings"
	"time"
)

var jsflowParams = map[string]string{
	"scale": "1",
}

type jsflowRenderer struct {
}

func (j *jsflowRenderer) BlockDescriptor() string {

	return "jsflow"

}

func (j *jsflowRenderer) HTMLHeaders() string {

	return `
    <script src="http://flowchart.js.org/flowchart-latest.js"></script>
    `
}

func (j *jsflowRenderer) RenderToBuffer(data string, params string) (string, error) {

	p := params2map(params, jsflowParams)

	data = strings.Replace(data, "\n", "\\n", -1)

	divID := fmt.Sprintf("div%v", time.Now().UnixNano())

	div := `
    <div id="` + divID + `"></div>
    <script>
        var diagram = flowchart.parse("` + data + `");
        diagram.drawSVG('` + divID + `');

        scale = ` + p["scale"] + `;
        svgNode = $($("#` + divID + `").children()[0]);
        svgNode.html("<g transform='scale("+scale+")'>"+svgNode.html()+"</g>");
        svgNode.attr("height", svgNode.attr("height") * scale );
        svgNode.attr("width", svgNode.attr("width") * scale );
        </script>
    `

	return div, nil

}
