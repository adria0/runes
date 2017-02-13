package render

import (
	"fmt"
	"time"
)

var jsseqParams = map[string]string{
	"scale": "1",
	"theme": "simple",
}

type jsseqRenderer struct {
}

func (j *jsseqRenderer) BlockDescriptor() string {

	return "jsseq"

}

func (j *jsseqRenderer) HTMLHeaders() string {

	return `
    <script src="https://cdnjs.cloudflare.com/ajax/libs/raphael/2.2.7/raphael.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.8.3/underscore-min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/js-sequence-diagrams/1.0.6/sequence-diagram-min.js"></script>
    `
}

func (j *jsseqRenderer) RenderToBuffer(data string, params string) (string, error) {

	p := params2map(params, jsseqParams)

	divID := fmt.Sprintf("div%v", time.Now().UnixNano())

	div := `
    <div id="` + divID + `">` + data + `</div>
    <script>
    $("#` + divID + `").sequenceDiagram({theme: '` + p["theme"] + `'});

    scale = ` + p["scale"] + `;
    svgNode = $($("#` + divID + `").children()[0]);
    svgNode.html("<g transform='scale("+scale+")'>"+svgNode.html()+"</g>");
    svgNode.attr("height", svgNode.attr("height") * scale );
    svgNode.attr("width", svgNode.attr("width") * scale );
    </script>
    `

	return div, nil

}
