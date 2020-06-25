package postman

import "bytes"

const (
	ENV_EXT      = "postman_environment.json"
	COLL_EXT     = "postman_collection.json"
	SCHEMA_2_1_0 = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)

type Header struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Warning  string `json:"warning"`
	Disabled bool   `json:"disabled"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Url struct {
	Raw   string     `json:"raw"`
	Host  []string   `json:"host"`
	Path  []string   `json:"path"`
	Query []KeyValue `json:"query"`
}

type Request struct {
	Method    string     `json:"method"`
	Header    []Header   `json:"header"`
	Url       Url        `json:"url"`
	Body      *Body      `json:"body"`
	Variables []KeyValue `json:"variable"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type ProtocolProfileBehavior struct {
	DisableBodyPruning bool `json:"disableBodyPruning"`
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Script struct {
	Id    string   `json:"id"`
	Lines []string `json:"exec"`
	Type  string   `json:"type"`
}

func (s *Script) String() string {
	var buf bytes.Buffer
	for _, l := range s.Lines {
		buf.WriteString(l)
		buf.WriteString("\n")
	}
	return buf.String()
}
