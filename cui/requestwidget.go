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
	r := rw.collection.Request
	fmt.Fprintln(v, r.Method, rw.collection.Name)
	fmt.Fprintf(v, "%s\n", r.Url.Raw)
	fmt.Fprintln(v, colorWhite, "------------------------------------", colorReset)
	for _, h := range r.Header {
		fmt.Fprintf(v, "\t%s - %s\n", h.Key, h.Value)
	}
	fmt.Fprintln(v, colorWhite, "------------------------------------", colorReset)
	if r.Body != nil {
		fmt.Fprintln(v, r.Body.Raw)
	}
	fmt.Fprintln(v, colorYellow, "---------- Script ------------------", colorReset)
	fmt.Fprint(v, colorBlue)
	for _, ev := range rw.collection.Events {
		for _, sl := range ev.Script.Lines {
			fmt.Fprintln(v, sl)
		}
	}
	fmt.Fprint(v, colorReset)
}
