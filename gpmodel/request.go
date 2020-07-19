package gpmodel

type RequestSpec struct {
	Name       string
	Method     string
	UrlPattern string
	Headers    []Header
	Body       string
}

type Header struct {
	Key, Value string
}
