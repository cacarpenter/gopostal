package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/cacarpenter/gopostal/util"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

type RequestWidget struct {
	request *gpmodel.RequestSpec
	vars map[string]string
}

func (rw *RequestWidget) Layout(v *gocui.View) {
	v.Clear()
	if rw.request == nil {
		return
	}

	r := rw.request
	fmt.Fprintln(v, colorCyan, r.Method, colorReset, r.Name)
	fmt.Fprintf(v, "%s\n", r.UrlPattern)
	fmt.Fprintf(v, "%s%s%s\n", colorPurple, util.ReplaceVariables(r.UrlPattern, rw.vars), colorReset)

	headers := tablewriter.NewWriter(v)
	headers.SetHeader([]string{"Key", "Value", "Send Value"})

	for _, h := range r.Headers {
		headers.Append([]string{h.Key, h.Value, util.ReplaceVariables(h.Value, rw.vars)})
	}
	headers.Render()
	fmt.Fprintln(v, r.Body)
	fmt.Fprintln(v, colorYellow, "---------- Script ------------------", colorReset)
	fmt.Fprint(v, colorBlue)
	fmt.Fprintln(v, r.PostScript.Text)
	fmt.Fprintln(v, colorReset)
}

func (rw *RequestWidget) UpdateVars() {

}