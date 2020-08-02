package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/cacarpenter/gopostal/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (app *GoPostal) CallRequest(reqSpec *gpmodel.RequestSpec, writer io.Writer) (*string, error) {
	app.logger.Println("Using Request Pattern", reqSpec.UrlPattern)
	interUrl := util.ReplaceVariables(reqSpec.UrlPattern, app.session.variables)
	app.logger.Println("calling", interUrl)
	httpClient := http.Client{}
	httpReq, err := http.NewRequest(reqSpec.Method, interUrl, strings.NewReader(reqSpec.Body))
	if err != nil {
		app.logger.Println("Bad Req")
		return nil, err
	}
	for _, hdr := range reqSpec.Headers {
		headerVal := util.ReplaceVariables(hdr.Value, app.session.variables)
		headerKey := hdr.Key
		// TODO this is something postman specific? Need this to be the Authorization Header with a Bearer token
		if headerKey == "Bearer" {
			headerKey = "Authorization"
			headerVal = "Bearer " + headerVal
		}
		httpReq.Header.Add(hdr.Key, headerVal)
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
		return nil, fmt.Errorf("got %d as response", httpResp.StatusCode)
	}

	strRcvBody := string(rcvBody)
	fmt.Fprintln(writer, strRcvBody)

	return &strRcvBody, nil
}
