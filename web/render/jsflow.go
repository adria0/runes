package render

import (
    "time"
    "fmt"
    "strings"
)

type jsflowRenderer struct {

}

func (j *jsflowRenderer) BlockDescriptor() string {

    return "jsflow"

}

func (j *jsflowRenderer) HtmlHeaders() string {

    return `
    <script src="http://github.com/DmitryBaranovskiy/raphael/raw/master/raphael-min.js"></script>
    <script src="http://flowchart.js.org/flowchart-latest.js"></script>
    `
}

func (j *jsflowRenderer) RenderToBuffer(data string, params string) (string, error) {

    data = strings.Replace(data,"\n","\\n",-1)

    divId := fmt.Sprintf("div%v",time.Now().UnixNano())

    div := `
    <div id="`+divId+`"></div>
    <script>
        var diagram = flowchart.parse("`+data+`");
        diagram.drawSVG('`+divId+`');
    </script>
    `

    return div, nil

}


