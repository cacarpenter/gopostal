package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (app *GoPostal) CallRequest(pmreq *postman.Request, writer io.Writer) (*string, error) {
	interUrl := util.ReplaceVariables(pmreq.Url.Raw, app.session.variables)
	app.logger.Println("calling", interUrl)
	httpClient := http.Client{}
	var sendBody io.Reader
	if pmreq.Body != nil {
		sendBody = strings.NewReader(pmreq.Body.Raw)
	}
	httpReq, err := http.NewRequest(pmreq.Method, interUrl, sendBody)
	if err != nil {
		app.logger.Println("Bad Req")
		return nil, err
	}
	for _, pmHeader := range pmreq.Header {
		headerVal := util.ReplaceVariables(pmHeader.Value, app.session.variables)
		headerKey := pmHeader.Key
		// TODO this is something postman specific? We need this to be the Authorization Header with a Bearer token
		if headerKey == "Bearer" {
			headerKey = "Authorization"
			headerVal = "Bearer "+headerVal
		}
		httpReq.Header.Add(pmHeader.Key, headerVal)
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		app.logger.Println(writer, "Bad Resp")
		return nil, err
	}
	defer httpResp.Body.Close()

	rcvBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Println("Error Reading Bad Read")
		return nil, err
	}

	fmt.Fprintln(writer, httpResp.StatusCode)
	// process events - need some rules around when to do this but for now just look for some basic valid responses
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 201 {
		// TODO return sendBody
		return nil, fmt.Errorf("Got %d as response", httpResp.StatusCode)
	}

	strRcvBody := string(rcvBody)
	fmt.Fprintln(writer, strRcvBody)

	return &strRcvBody, nil
}
