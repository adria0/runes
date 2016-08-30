package render

import (
    "time"
    "fmt"
)

type jsseqRenderer struct {

}

func (j *jsseqRenderer) BlockDescriptor() string {

    return "jsseq"

}

func (j *jsseqRenderer) HtmlHeaders() string {

    return `
    <script src="https://bramp.github.io/js-sequence-diagrams/js/raphael-min.js"></script>
    <script src="https://bramp.github.io/js-sequence-diagrams/js/underscore-min.js"></script>
    <script src="https://bramp.github.io/js-sequence-diagrams/js/sequence-diagram-min.js"></script>
    `
}

func (j *jsseqRenderer) RenderToBuffer(data string, params string) (string, error) {

    divId := fmt.Sprintf("div%v",time.Now().UnixNano())

    div := `
    <div id="`+divId+`">` + data + `</div>
    <script>
    $("#`+divId+`").sequenceDiagram({theme: 'simple'});
    </script>
    `

    return div, nil

}


