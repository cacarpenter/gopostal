package cui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ResponseWidget struct {
	requestUrl string
	calling    bool
	statusCode int
	body       string
}

func (rw *ResponseWidget) SetRequesting(url string) {
	rw.requestUrl = url
	rw.calling = true
}

func (rw *ResponseWidget) SetResponse(status int, body string) {
	rw.calling = false
	rw.body = body
	rw.statusCode = status
}

func (rw *ResponseWidget) Layout(v *gocui.View) {
	if rw.calling {
		fmt.Fprintf(v, "Calling %q\n", rw.requestUrl)
	} else {
		fmt.Fprintf(v, "%d\n%s\n", rw.statusCode, rw.body)
	}
}
