package postman

const SCHEMA_2_1_0 = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

type Collection struct {
	Info  CollectionInfo   `json:"info"`
	Items []CollectionItem `json:"item"`
}

type CollectionInfo struct {
	PostmanId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type CollectionItem struct {
	Name     string           `json:"name"`
	Children []CollectionItem `json:"item"`
	Events   []Event          `json:"event"`
	Request  *Request         `json:"request"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Query struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Url struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []Query  `json:"query"`
}

type Request struct {
	Method string   `json:"method"`
	Header []Header `json:"header"`
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
}

type Script struct {
	Id    string   `json:"id"`
	Lines []string `json:"exec"`
	Type  string   `json:"type"`
}
