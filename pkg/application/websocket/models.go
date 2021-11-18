package websocket

type shortURLDataIn struct {
	RequestType string `json:"request_type"`
	Request     struct {
		URLs []string `json:"urls"`
	} `json:"request"`
}

type shortURLDataOut struct {
	ResponseType string `json:"response_type"`
	Response     struct {
		URLs []string `json:"urls"`
	} `json:"response"`
}
