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
	rs.Headers = make([]gpmodel.Header, len(pmReq.Header))
	for i, pmHeader := range pmReq.Header {
		rs.Headers[i] = gpmodel.Header{Key: pmHeader.Key, Value: pmHeader.Value}
	}
	if pmReq.Body != nil {
		pmBod := pmReq.Body
		rs.Body = pmBod.Raw
		if pmBod.Options != nil && pmBod.Options.Raw.Language == "json"{
			rs.Headers = append(rs.Headers, gpmodel.Header{Key: "Content-Type", Value: "application/json"})
		}
	}
	return rs
}
