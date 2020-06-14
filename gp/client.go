package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
)

func CallRequest(req *postman.Request) error {
	fmt.Println("call ", req.Url.Raw)

	return nil
}

func replaceVariables(raw string) string {
	sess := CurrentSession()

}
