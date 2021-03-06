package gpmodel

type RequestSpec struct {
	Name       string
	Method     string
	UrlPattern string
	Headers    []Header
	Body       string
	PostScript Script
}

type Header struct {
	Key, Value string
}

type Script struct {
	Type string
	Text string
}
