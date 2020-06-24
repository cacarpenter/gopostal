package cui

import (
	"github.com/cacarpenter/gopostal/gp"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

type VariablesWidget struct {
}

func (vw *VariablesWidget) Layout(view *gocui.View) {
	sess := gp.CurrentSession()

	table := tablewriter.NewWriter(view)
	table.SetHeader([]string{"Name", "Value"})
	table.SetBorder(false)

	for _, v := range sess.Array() {
		table.Append(v)
	}
	table.Render()
}
