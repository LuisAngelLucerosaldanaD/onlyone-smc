package viewer

type RequestViewer struct {
	Url string `json:"url"`
	Ttl int    `json:"ttl"`
}

type ResponseViewer struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}
