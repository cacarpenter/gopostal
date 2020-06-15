package gp

import (
	"bytes"
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"net/http"
	"regexp"
	"strings"
)

func CallRequest(pmreq *postman.Request) error {
	fmt.Println("call", replaceVariables(pmreq.Url.Raw))

	interUrl := replaceVariables(pmreq.Url.Raw)
	httpClient := http.Client{}
	httpReq, err := http.NewRequest(pmreq.Method, interUrl, strings.NewReader(pmreq.Body.Raw))
	if err != nil {
		return err
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}

	fmt.Println(httpResp.StatusCode)
	return nil
}

func replaceVariables(raw string) string {
	r, err := regexp.Compile("\\{\\{\\w+\\}\\}")
	if err != nil {
		panic(err)
	}

	sess := CurrentSession()

	var buffer bytes.Buffer
	fasi := r.FindAllStringIndex(raw, -1)
	prevIdx := 0
	for _, matchIndices := range fasi {
		// TODO should be able to remove this
		if len(matchIndices) != 2 {
			panic(fmt.Errorf("Thought length would be 2 but it was %d\n", len(matchIndices)))
		}

		// chars before this match
		if matchIndices[0] > prevIdx {
			buffer.WriteString(raw[prevIdx:matchIndices[0]])
		}
		varstr := raw[matchIndices[0]+2 : matchIndices[1]-2]
		val := sess.Get(varstr)
		buffer.WriteString(val)

		prevIdx = matchIndices[1]
	}
	// some literal text remains
	if prevIdx < len(raw)-1 {
		buffer.WriteString(raw[prevIdx:len(raw)])
	}

	return buffer.String()
}
