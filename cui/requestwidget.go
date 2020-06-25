package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gp"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

type RequestWidget struct {
	collection *postman.Collection
}

func (rw *RequestWidget) Layout(v *gocui.View) {
	v.Clear()
	if rw.collection == nil || rw.collection.Request == nil {
		return
	}

	currVals := gp.CurrentSession().Map()

	r := rw.collection.Request
	fmt.Fprintln(v, colorCyan, r.Method, colorReset, rw.collection.Name)
	fmt.Fprintf(v, "%s\n", r.Url.Raw)
	fmt.Fprintf(v, "%s%s%s\n", colorPurple, util.ReplaceVariables(r.Url.Raw, currVals), colorReset)

	headers := tablewriter.NewWriter(v)
	headers.SetHeader([]string{"Key", "Value", "Send Value"})

	for _, h := range r.Header {
		headers.Append([]string{h.Key, h.Value, util.ReplaceVariables(h.Value, currVals)})
	}
	headers.Render()
	if r.Body != nil {
		fmt.Fprintln(v, r.Body.Raw)
	}
	fmt.Fprintln(v, colorYellow, "---------- Script ------------------", colorReset)
	fmt.Fprint(v, colorBlue)
	for _, ev := range rw.collection.Events {
		for _, script := range ev.Script.Lines {
			fmt.Fprintln(v, script)
		}
	}
	fmt.Fprintln(v, colorReset)
}
