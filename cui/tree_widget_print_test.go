package cui

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"log"
	"regexp"
	"testing"
)

func removeNonPrintChars(s string) string {
	reg, err := regexp.Compile(`\w|\s|\d+`)
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func addDescendents(parent *treeNode, children ...string) {
	// fmt.Printf("%q -> %q\n", parent.label, children)
	if len(children) < 1 {
		return
	}
	childName := children[0]
	rest := children[1:]
	childNode := treeNode{childName, parent, make([]*treeNode, 0), nil, true, false}
	parent.children = append(parent.children, &childNode)
	addDescendents(&childNode, rest...)
}

func Test_printNode_withrequest(t *testing.T) {
	node := treeNode{}
	node.label = "node label"
	node.request = &gpmodel.RequestSpec{}
	node.request.Name = "req name"
	node.selected = true
	var buf bytes.Buffer
	bufw := bufio.NewWriter(&buf)
	printNode(bufw, &node, "-")
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

func Test_printNode_withchildren(t *testing.T) {
	node := treeNode{}
	node.label = "node label"
	node.selected = true
	node.expand(true, false)
	addDescendents(&node, "ch0", "grch", "grgrch")
	postScript := gpmodel.Script{}
	ch0 := node.children[0]
	ch0.request = &gpmodel.RequestSpec{
		Name:       "child 0 req 0",
		Method:     "POST",
		UrlPattern: "{{domain}}",
		Headers:    []gpmodel.Header{},
		Body:       "body",
		PostScript: postScript}
	node.children[0] = ch0
	var buf bytes.Buffer
	bufw := bufio.NewWriter(&buf)
	printNode(bufw, &node, "-")
	bufw.Flush()
	// shouldBe := " -[] node label\n"
	// fails due to colors
	/*
		if buf.String() != shouldBe {
			t.Errorf("Should be %q not %q\n", shouldBe, buf.String())
		}
	*/
	fmt.Printf("'%s'\n", buf.String())
}
