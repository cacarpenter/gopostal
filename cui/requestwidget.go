package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
)

type RequestWidget struct {
	collection *postman.Collection
}

func (rw *RequestWidget) Layout(v *gocui.View) {
	v.Clear()
	if rw.collection == nil || rw.collection.Request == nil {
		return
	}
	fmt.Fprintln(v, rw.collection.Name)
	r := rw.collection.Request
	fmt.Fprintf(v, "%s - %s\n", r.Method, r.Url.Raw)
	if r.Body != nil {
		fmt.Fprintln(v, r.Body.Raw)
	}
}
