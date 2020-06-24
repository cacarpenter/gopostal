package gp

import (
	"bytes"
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func CallRequest(pmreq *postman.Request, writer io.Writer) (*string, error) {
	interUrl := replaceVariables(pmreq.Url.Raw)
	fmt.Fprintln(writer, "call", interUrl)
	httpClient := http.Client{}
	var sendBody io.Reader
	if pmreq.Body != nil {
		sendBody = strings.NewReader(pmreq.Body.Raw)
	}
	httpReq, err := http.NewRequest(pmreq.Method, interUrl, sendBody)
	if err != nil {
		fmt.Fprintln(writer, "Bad Req")
		return nil, err
	}
	for _, pmHeader := range pmreq.Header {
		httpReq.Header.Add(pmHeader.Key, pmHeader.Value)
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Fprintln(writer, "Bad Resp")
		return nil, err
	}
	defer httpResp.Body.Close()

	rcvBody, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		fmt.Fprintln(writer, "Bad Read")
		return nil, err
	}

	fmt.Fprintln(writer, httpResp.StatusCode)
	// process events - need some rules around when to do this but for now just look for some basic valid responses
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 201 {
		// TODO return sendBody
		return nil, fmt.Errorf("Got %d as response", httpResp.StatusCode)
	}


	strRcvBody := string(rcvBody)
	fmt.Println(writer, strRcvBody)

	return &strRcvBody, nil
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
