package cui

import (
	"github.com/cacarpenter/gopostal/util"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

type VariablesWidget struct {
	vars               [][]string
	editMode           bool
	selectedRow        int
	selectedValueOrKey bool
}

func (vw *VariablesWidget) Layout(view *gocui.View) {
	table := tablewriter.NewWriter(view)
	table.SetHeader([]string{"Name", "Value"})
	table.SetBorder(false)

	if vw.editMode {
		rows := make([]string, 1)
		rows[0] = "R1"
		colors := make([]tablewriter.Colors, 1)
		colors[0] = tablewriter.Colors{}

		table.Rich(rows, colors)
	} else {
		table.AppendBulk(vw.vars)
	}

	table.Render()
}

func (vw *VariablesWidget) SetVariables(vars [][]string) {
	vw.vars = util.SortArray(vars)
}

func (vw *VariablesWidget) SetVariable(k, v string) {
	for _, vk := range vw.vars {
		if vk[0] == k {
			vk[1] = v
			break
		}
	}
}
