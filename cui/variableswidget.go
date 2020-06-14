package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gp"
	"github.com/jroimartin/gocui"
)

type VariablesWidget struct {
}

func (vw *VariablesWidget) Layout(view *gocui.View) {
	sess := gp.CurrentSession()
	fmt.Fprintln(view, len(sess.Vars()))
	for key, val := range sess.Vars() {
		fmt.Fprintln(view, key, val)
	}
}
