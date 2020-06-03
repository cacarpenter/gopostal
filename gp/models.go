package gp

type CollectionInfo struct {
	PostmanId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type Header struct {
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

type Item struct {
	Name string `json:"name"`
}

type PostmanCollection struct {
	Info CollectionInfo `json:"info"`
}
