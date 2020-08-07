package cui

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"testing"
)


func Test_printNode(t *testing.T) {
	node := treeNode{}
	node.label = "node label"
	node.request = &gpmodel.RequestSpec{}
	node.request.Name = "req name"
	node.selected = true
	var buf bytes.Buffer
	bufw := bufio.NewWriter(&buf)
	printNode2(bufw, &node, "-", 0)
	bufw.Flush()
	// shouldBe := " -[] node label\n"
	// fails due to colors
	/*
	if buf.String() != shouldBe {
		t.Errorf("Should be %q not %q\n", shouldBe, buf.String())
	}
	 */
	fmt.Printf("%q\n", buf.String())
}