package postman

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
)

func printCollectionItem(pad string, ci *Collection) {
	fmt.Println(pad+"* Name :", ci.Name)
	if ci.Request != nil {
		fmt.Println(pad+"* Req  : ", ci.Request)
	}
	for _, child := range ci.Items {
		printCollectionItem(pad+pad, child)
	}
}

func Print(coll *Collection) {
	fmt.Println("Name:", coll.Info.Name)
	for _, i := range coll.Items {
		printCollectionItem("  ", i)
	}
}

func NewRequestSpec(pmReq *Request) *gpmodel.RequestSpec {
	rs := new(gpmodel.RequestSpec)
	rs.Method = pmReq.Method
	rs.UrlPattern = pmReq.Url.Raw
	if pmReq.Body != nil {
		rs.Body = pmReq.Body.Raw
	}
	return rs
}
